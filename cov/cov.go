package cov

import (
	"bcov/bam"
	"bcov/db"
	"bcov/utils"
	"fmt"
	"io"
	"log"

	"github.com/biogo/hts/sam"
	"github.com/fatih/color"
)

func bamWorker(bamReader *bam.Reader, regions []Region, rChan chan<- Region) {
	spinner := utils.NewSpinner(fmt.Sprintf("reading %s", bamReader.Path))
	spinner.Start()
	defer spinner.StopDuration()

	for len(regions) != 0 {
		rec, err := bamReader.Read()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatalf("error reading bam: %v", err)
		}

		readChromosome := rec.Ref.Name()
		readStart, readEnd := uint64(rec.Pos+1), uint64(rec.End())
		if readStart%1000 == 0 {
			go spinner.Message(fmt.Sprintf("chromosome %s pos %d Mpb (%.2f%%)", readChromosome, readStart/100000, 100*float64(readStart)/float64(utils.CHROMOSOME_LENGTHS[utils.ChromosomeIndex(readChromosome)-1])))
		}

		// Sequencing read SAM flags to exclude from counting, bitwise
		flags := !(rec.Flags&sam.Duplicate == sam.Duplicate)

		if flags {
			if utils.ChromosomeIndex(regions[0].Chromosome) > utils.ChromosomeIndex(readChromosome) {
				continue
			}

			// Add 1 to positions that fall within this read
			// (NOTE: iterating regions because >>a read can fall inside in more than 1 region<<)
			for _, r := range regions {
				overlap := r.Chromosome == readChromosome && readStart <= r.End && readEnd >= r.Start
				if overlap {
					r.AddDepthFromTo(Position(readStart), Position(readEnd), 1)
				}
				if !overlap && (r.Start > readEnd || r.Chromosome != readChromosome) {
					break
				}
			}
		}

		// Check if first region is way past the current read coordinates (we are
		// done calculating this region positions depth) to delete it from memory and store in DB
		//
		// sameChromosome := readChromosome == regions[0].Region.Chromosome
		// pastPosition := readStart > regions[0].Region.End
		// pastChromosome := utils.CHROMOSOME_INDEX[readChromosome] > utils.CHROMOSOME_INDEX[regions[0].Region.Chromosome]
		// if (sameChromosome && pastPosition) || pastChromosome {

		if (readChromosome == regions[0].Chromosome && readStart > regions[0].End) || utils.ChromosomeIndex(readChromosome) > utils.ChromosomeIndex(regions[0].Chromosome) {
			rChan <- regions[0]
			regions = regions[1:]
		}
	}

	close(rChan)
}

func GetRegionDepth(bamPath string, chromosome string, start uint64, end uint64) Region {
	bamReader, err := bam.NewReader(bamPath)
	if err != nil {
		log.Fatal(err)
	}

	region := NewRegion(chromosome, start, end)
	rChan := make(chan Region)
	go bamWorker(bamReader, []Region{region}, rChan)

	for range rChan {
		return region
	}

	return region
}

func Load(bamPath string) {
	bamReader, err := bam.NewReader(bamPath)
	if err != nil {
		log.Fatal(err)
	}

	hash := bamReader.SHA256sum()
	bamFileID, created := db.StoreFile(bamReader.Filename, hash)
	if !created {
		fmt.Printf("File %s (%s) already exists in database\n", bamReader.Path, hash)
		fmt.Println()
		fmt.Printf("If you want to load this file again, before run -delete-bam %s\n", hash)
		return
	}
	regions := NewRegionsFromDB()
	rChan := make(chan Region)
	go bamWorker(bamReader, regions, rChan)

	for r := range rChan {
		r.StoreDepthCoverages(bamFileID)
	}
}

func prettyPrint(region db.Region, depthCoverages map[int]float64) {
	fmt.Println()
	fmt.Print("Region ")
	color.Red("%s:%d-%d (%d)", region.Chromosome, region.Start, region.End, region.GeneID)
	fmt.Println(depthCoverages)
	fmt.Println()
}
