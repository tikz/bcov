package cov

import (
	"bcov/bam"
	"bcov/db"
	"bcov/utils"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/biogo/hts/sam"
	"github.com/fatih/color"
)

type Position uint64
type Depth uint64

type RegionDepths struct {
	Region        db.Region
	PositionDepth map[Position]Depth
}

type Loader struct {
	BAM struct {
		Path   string
		Reader *bam.Reader
	}
}

// NewLoader returns a new Loader struct.
func NewLoader(bamPath string, bedPath string) *Loader {
	bl := &Loader{}

	bamReader, err := bam.NewReader(bamPath)
	if err != nil {
		log.Fatal(err)
	}
	bl.BAM.Path = bamPath
	bl.BAM.Reader = bamReader

	return bl
}

// nextRegion advances the BED reader and returns the next region from the BED file.
// Included in the returned RegionDepths struct are the parsed coordinates and a map of
// positions to depths, initialized with zeroes.
// If the BED file is exhausted, io.EOF is returned.
// func (l *Loader) nextRegion() (RegionDepths, error) {
// 	region, err := l.BED.Reader.Read()
// 	if err == io.EOF {
// 		return RegionDepths{}, err
// 	}

// 	positionDepths := make(map[Position]Depth)
// 	for i := region.Start; i <= region.End; i++ {
// 		positionDepths[Position(i)] = 0
// 	}

// 	return RegionDepths{Region: &region, PositionDepth: positionDepths}, nil
// }

// func Start() {

// 	bamReader, err := bam.NewReader("B0434109PTC_Result/B0434109PTC.bam")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	rdJobsCh := make(chan RegionDepths)
// 	rdResultsCh := make(chan RegionDepths)
// 	bamReadsCh := make(chan *sam.Record)
// 	regions := db.GetRegions()

// 	// Launch workers
// 	go bamWorker(rdJobsCh, bamReadsCh, rdResultsCh)
// 	go dbStorer(bamReader.Filename, bamReader.SHA256, rdResultsCh)

// 	go func() {
// 		for _, region := range regions {
// 			positionDepths := make(map[Position]Depth)
// 			for i := region.Start; i <= region.End; i++ {
// 				positionDepths[Position(i)] = 0
// 			}
// 			rdJobsCh <- RegionDepths{Region: region, PositionDepth: positionDepths}
// 		}
// 	}()

// 	// Feed the channel with BAM reads
// 	for {
// 		read, err := bamReader.Read()
// 		if err == io.EOF {
// 			close(rdJobsCh)
// 			close(bamReadsCh)
// 			break
// 		}

// 		bamReadsCh <- read
// 	}
// }

func dbStorer(bamFilename string, bamHash string, regionDepths <-chan RegionDepths) {
	bamFileID := db.StoreFile(bamFilename, bamHash)
	for rd := range regionDepths {
		depthCoverages := make(map[int]float64)
		for i := 1; i <= 100; i++ {
			count := 0

			for _, depth := range rd.PositionDepth {
				if int(depth) >= i {
					count++
				}
			}
			depthCoverages[i] = float64(count) / float64(rd.Region.End-rd.Region.Start+1)
		}
		db.StoreDepthCoverages(bamFileID, rd.Region.ID, depthCoverages)
	}
}

func pof() {

}

func bamWorker(rdIn <-chan RegionDepths, bamReadsCh <-chan *sam.Record, rdOut chan<- RegionDepths) {
	rd := <-rdIn

	spinner := utils.NewSpinner(fmt.Sprintf("reading"))
	spinner.Start()
	defer spinner.StopDuration()

	for {
		rec := <-bamReadsCh

		chromosome := rec.Ref.Name()
		readStart, readEnd := uint64(rec.Pos), uint64(rec.End())
		if readStart%1000 == 0 {
			spinner.Message(fmt.Sprintf("chromosome %s pos %d", chromosome, readStart))
		}

		// Read falls outside of current region
		if readStart > rd.Region.End && chromosome == rd.Region.Chromosome {
			if rd.Region.ID == 209414 {
				fmt.Println(rd.Region)
				fmt.Println(rd.PositionDepth)
				fmt.Println()
				os.Exit(0)
			}
			rdOut <- rd
			rd = <-rdIn
		}

		// Sequencing read SAM flags to exclude from counting, bitwise
		flags := !(rec.Flags&sam.Duplicate == sam.Duplicate)

		// Add 1 to positions that fall within this read
		overlap := rd.Region.Chromosome == chromosome && readStart <= rd.Region.End && readEnd >= rd.Region.Start
		if flags && overlap {
			for pos := range rd.PositionDepth {
				if readStart <= uint64(pos) && readEnd >= uint64(pos) {
					rd.PositionDepth[pos]++
				}
			}
		}

		// if readStart > 35110 {
		// 	fmt.Println(rd.Region)
		// 	fmt.Println(rd)
		// 	os.Exit(0)
		// }
	}
}

func max(a uint64, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

func min(a uint64, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}

func initializeRegionDepths() []RegionDepths {
	regions := db.GetRegions()
	var regionDepths []RegionDepths
	for _, region := range regions {
		regionDepths = append(regionDepths, RegionDepths{Region: region, PositionDepth: make(map[Position]Depth)})
	}

	return regionDepths
}

func Load(bamPath string) {
	regions := initializeRegionDepths()
	// currentRegion, regions := regions[0], regions[1:]
	// currentRD := initializeRegionDepths()

	bamReader, err := bam.NewReader(bamPath)
	if err != nil {
		log.Fatal(err)
	}

	spinner := utils.NewSpinner(fmt.Sprintf("reading %s", bamPath))
	spinner.Start()
	defer spinner.StopDuration()
	for {
		if len(regions) == 0 {
			break
		}

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
			spinner.Message(fmt.Sprintf("chromosome %s pos %d", readChromosome, readStart))
		}

		// Sequencing read SAM flags to exclude from counting, bitwise
		flags := !(rec.Flags&sam.Duplicate == sam.Duplicate)

		if flags {
			// Add 1 to positions that fall within this read
			// (NOTE: iterating regions because >>a read can fall inside in more than 1 region<<)
			for _, r := range regions {
				overlap := r.Region.Chromosome == readChromosome && readStart <= r.Region.End && readEnd >= r.Region.Start
				if overlap {
					start, end := max(r.Region.Start, readStart), min(r.Region.End, readEnd)
					for i := start; i <= end; i++ {
						if _, ok := r.PositionDepth[Position(i)]; !ok {
							r.PositionDepth[Position(i)] = 0
						}
						r.PositionDepth[Position(i)]++
					}
				}
				if !overlap && (r.Region.Start > readEnd || r.Region.Chromosome != readChromosome) {
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

		if (readChromosome == regions[0].Region.Chromosome && readStart > regions[0].Region.End) || utils.ChromosomeIndex(readChromosome) > utils.ChromosomeIndex(regions[0].Region.Chromosome) {
			fmt.Println(regions[0].PositionDepth)
			regions = regions[1:]
		}

	}

}

func prettyPrint(region db.Region, depthCoverages map[int]float64) {
	fmt.Println()
	fmt.Print("Region ")
	color.Red("%s:%d-%d (%d)", region.Chromosome, region.Start, region.End, region.GeneID)
	fmt.Println(depthCoverages)
	fmt.Println()
}
