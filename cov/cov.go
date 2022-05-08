package cov

import (
	"bcov/bam"
	"bcov/db"
	"bcov/utils"
	"fmt"
	"io"
	"log"
	"strings"

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

		readChromosome := strings.ReplaceAll(rec.Ref.Name(), "chr", "")
		readStart, readEnd := uint64(rec.Pos+1), uint64(rec.End())
		if readStart%1000 == 0 {
			chromosomeIndex := utils.ChromosomeIndex(readChromosome)
			if chromosomeIndex != 99 {
				go spinner.Message(fmt.Sprintf("chromosome %s pos %d Mpb (%.2f%%)", readChromosome, readStart/100000, 100*float64(readStart)/float64(utils.CHROMOSOME_LENGTHS[chromosomeIndex-1])))
			}
		}

		// SAM flags to exclude
		flags := !(rec.Flags&sam.Duplicate == sam.Duplicate)

		if flags {
			if utils.ChromosomeIndex(regions[0].Chromosome) > utils.ChromosomeIndex(readChromosome) {
				continue
			}

			// Regions positions that fall within this read
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

		// Check if current region is way past the current read coordinates (done counting)
		//
		// 1. sameChromosome := readChromosome == regions[0].Region.Chromosome
		// 2. pastRegion := readStart > regions[0].Region.End
		// 3. pastChromosome := utils.CHROMOSOME_INDEX[readChromosome] > utils.CHROMOSOME_INDEX[regions[0].Region.Chromosome]

		// Condition: if (sameChromosome && pastRegion) || pastChromosome {
		//
		// (not defined as actual variables for performance, the evaluation can be resolved prematurely)

		if (readChromosome == regions[0].Chromosome && readStart > regions[0].End) || utils.ChromosomeIndex(readChromosome) > utils.ChromosomeIndex(regions[0].Chromosome) {
			rChan <- regions[0]
			regions = regions[1:]
		}
	}
	spinner.StopDuration()
	close(rChan)
}

func Load(bamPath string, kit string) {
	bamReader, err := bam.NewReader(bamPath)
	if err != nil {
		log.Fatal(err)
	}

	hash := bamReader.SHA256sum()

	kitID, _ := db.StoreKit(kit)
	bamFileID, created := db.StoreFile(bamReader.Filename, hash, kitID, bamReader.Size)
	if !created {
		fmt.Printf("File %s (%s) already exists in database\n", bamReader.Path, hash)
		fmt.Println()
		fmt.Printf("If you want to load this file again, run -delete-bam %s\n", hash)
		return
	}
	regions := NewRegionsFromDB()
	rChan := make(chan Region)
	go bamWorker(bamReader, regions, rChan)

	for r := range rChan {
		r.StoreDepthCoverages(bamFileID)
		r.StoreReadCounts(bamFileID)
	}
}

func RegionDepth(bamPath string, chromosome string, start uint64, end uint64) Region {
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

func prettyPrint(region db.Region, depthCoverages map[int]float64) {
	fmt.Println()
	fmt.Print("Region ")
	color.Red("%s:%d-%d (%d)", region.Chromosome, region.Start, region.End, region.GeneID)
	fmt.Println(depthCoverages)
	fmt.Println()
}
