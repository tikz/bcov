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

// bamWorker handles the parsing of a single BAM file, given a slice of exons to look for.
// The exon slice is >expected< to be sorted by karyotypic order (see utils/chromosome.go ChromosomeIndex() function)
func bamWorker(bamReader *bam.Reader, exons []Exon, rChan chan<- Exon) {
	spinner := utils.NewSpinner(fmt.Sprintf("reading %s", bamReader.Path))
	spinner.Start()
	defer spinner.StopDuration()

	for len(exons) != 0 {
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
			if utils.ChromosomeIndex(exons[0].Chromosome) > utils.ChromosomeIndex(readChromosome) {
				continue
			}

			// Exons positions that fall within this read
			for _, r := range exons {
				overlap := r.Chromosome == readChromosome && readStart <= r.End && readEnd >= r.Start
				if overlap {
					r.AddDepthFromTo(Position(readStart), Position(readEnd), 1)
				}
				if !overlap && (r.Start > readEnd || r.Chromosome != readChromosome) {
					break
				}
			}
		}

		// Logic and pseudocode:
		//
		// Check if the current exon is way past the current read coordinates (done counting)
		//
		// 1. sameChromosome := readChromosome == exons[0].Exon.Chromosome
		// 2. pastExon := readStart > exons[0].Exon.End
		// 3. pastChromosome := utils.CHROMOSOME_INDEX[readChromosome] > utils.CHROMOSOME_INDEX[exons[0].Exon.Chromosome]

		// Condition: if (sameChromosome && pastExon) || pastChromosome {
		//
		// (not defined as actual variables for performance, the evaluation can be resolved prematurely)

		if (readChromosome == exons[0].Chromosome && readStart > exons[0].End) || utils.ChromosomeIndex(readChromosome) > utils.ChromosomeIndex(exons[0].Chromosome) {
			rChan <- exons[0]
			exons = exons[1:]
		}
	}
	spinner.StopDuration()
	close(rChan)
}

// Load manages the loading and parsing of a BAM file given a path and a kit name.
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
	exons := NewExonsFromDB()
	rChan := make(chan Exon)
	go bamWorker(bamReader, exons, rChan)

	for r := range rChan {
		r.StoreDepthCoverages(bamFileID)
		r.StoreReadCounts(bamFileID)
	}
}

// ExonDepth calculates the read depth given BAM file and chromosomic coordinates, returning an Exon instance.
func ExonDepth(bamPath string, chromosome string, start uint64, end uint64) Exon {
	bamReader, err := bam.NewReader(bamPath)
	if err != nil {
		log.Fatal(err)
	}

	exon := NewExon(chromosome, start, end)
	rChan := make(chan Exon)
	go bamWorker(bamReader, []Exon{exon}, rChan)

	for range rChan {
		return exon
	}

	return exon
}

// helper to track progress
func prettyPrint(exon db.Exon, depthCoverages map[int]float64) {
	fmt.Println()
	fmt.Print("Exon ")
	color.Red("%s:%d-%d (%d)", exon.Chromosome, exon.Start, exon.End, exon.GeneID)
	fmt.Println(depthCoverages)
	fmt.Println()
}
