package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/biogo/hts/sam"
	"github.com/fatih/color"
)

var (
	testDb       = flag.Bool("test-db", false, "Test database connection")
	fetchRegions = flag.Bool("fetch-regions", false, "Fetch regions from Ensembl and store in DB")
	web          = flag.Bool("web", false, "Run web server")
	bam          = flag.String("bam", "", "load BAM file into the database")
	deleteBam    = flag.String("delete-bam", "", "delete BAM file from the database with SHA256")
	region       = flag.String("region", "", "Print per position depth of a given range, expressed as <chromosome>:<start>-<end>")
	help         = flag.Bool("help", false, "Display help")
)

func init() {
	d := color.New(color.FgRed, color.Bold)
	d.Println("Bcov")
	fmt.Println("Version: ", Version, "\t\tCommit: ", CommitHash)
	fmt.Println()
}

const maxFlag = int(^sam.Flags(0))

func main() {
	flag.Parse()
	if *help {
		fmt.Println("Usage:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *testDb {
		cliTestDB()
		os.Exit(0)
	}

	if *fetchRegions {
		cliFetchRegions()
		os.Exit(0)
	}

	if *region != "" {
		cliRegion()
		os.Exit(0)
	}

	if *bam != "" {
		cliLoadBAM()
		os.Exit(0)
	}

	if *deleteBam != "" {
		cliDeleteBam()
		os.Exit(0)
	}

	fmt.Println("Usage:")
	flag.PrintDefaults()
}

// func main2() {
// 	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
// 	if err != nil {
// 		panic("failed to connect database")
// 	}

// 	// Migrate the schema
// 	db.AutoMigrate(&BAMFile{}, &Gene{}, &Region{}, &DepthCoverage{})

// Create
// bamfile := BAMFile{Name: "test.bam", SHA256Sum: "test2"}
// db.Create(&bamfile)

// gene := Gene{Name: "CACA23"}
// db.Create(&gene)

// region := &Region{GeneID: gene.ID}
// db.Create(&region)

// depthCoverage := &DepthCoverage{BAMFileID: bamfile.ID, RegionID: region.ID, Depth: 1, Coverage: 1.0}
// db.Create(&depthCoverage)

// depthCoverage2 := &DepthCoverage{BAMFileID: bamfile.ID, RegionID: region.ID, Depth: 2, Coverage: 2.0}
// db.Create(&depthCoverage2)

// // Read
// var product Product
// db.First(&product, 1)                 // find product with integer primary key
// db.First(&product, "code = ?", "D42") // find product with code D42

// // Update - update product's price to 200
// db.Model(&product).Update("Price", 200)
// // Update - update multiple fields
// db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
// db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

// // Delete - delete product
// db.Delete(&product, 1)
// }

func mainx() {

	// cov.Load("B0434109PTC_Result/B0434109PTC.bam")
	// cov.GetRegionDepth("B0434109PTC_Result/B0434109PTC.bam", "4", 103910000, 103913000)

}

func main2() {
	// spinner := NewSpinner("test")
	// spinner2 := NewSpinner("test2")
	// spinner.Start()
	// spinner2.Start()
	// time.Sleep(time.Second * 5)
	// spinner.Stop()
	// spinner2.Stop()
	// reader, err := bed.NewReader("B0434109PTC_Result/test.bed")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(reader.Read())
	// fmt.Println(reader.Read())
	// fmt.Println(reader.Read())
	// l := cov.NewLoader("B0434109PTC_Result/B0434109PTC.bam", "B0434109PTC_Result/bigtest.bed")
	// l.Run()
}

// func main3() {

// 	flag.Parse()
// 	if *help {
// 		flag.Usage()
// 		os.Exit(0)
// 	}

// 	cfg := yacspin.Config{
// 		Frequency:       100 * time.Millisecond,
// 		CharSet:         yacspin.CharSets[59],
// 		Suffix:          " file " + *file,
// 		SuffixAutoColon: true,
// 		Message:         "opening",
// 		StopCharacter:   "✓",
// 		StopColors:      []string{"fgGreen"},
// 	}

// 	spinner, _ := yacspin.New(cfg)
// 	// handle the error
// 	spinner.Start()

// 	// spinner.Message("calculating file SHA256 sum")
// 	// fmt.Println(SHA256FileHash(*file))

// 	b, err := bam.NewReader(*file)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// var n int
// 	// var depth int

// 	regions, err := bed.LoadBED("B0434109PTC_Result/test.bed")
// 	if err != nil {
// 		panic(err)
// 	}

// 	type Position uint64
// 	type Depth uint64

// 	regionIndex := 0
// 	positionDepths := make(map[Position]Depth)
// 	for i := regions[regionIndex].Start; i <= regions[regionIndex].End; i++ {
// 		positionDepths[Position(i)] = 0
// 	}

// 	for {
// 		rec, err := b.Read()
// 		if err == io.EOF {
// 			break
// 		}

// 		chromosome := rec.Ref.Name()
// 		readStart, readEnd := uint64(rec.Pos), uint64(rec.End())

// 		if readStart%1000 == 0 {
// 			spinner.Message(fmt.Sprintf("chromosome %s pos %d", chromosome, readStart))
// 		}

// 		if readStart > regions[regionIndex].End && chromosome == regions[regionIndex].Chromosome {

// 			// fmt.Println(positionDepths)
// 			regionIndex++

// 			if regionIndex > len(regions)-1 {
// 				break
// 			}

// 			positionDepths = make(map[Position]Depth)
// 			for i := regions[regionIndex].Start; i <= regions[regionIndex].End; i++ {
// 				positionDepths[Position(i)] = 0
// 			}
// 		}

// 		flags := !(rec.Flags&sam.Duplicate == sam.Duplicate)
// 		overlap := regions[regionIndex].Chromosome == chromosome && readStart <= regions[regionIndex].End && readEnd >= regions[regionIndex].Start
// 		if flags && overlap {
// 			for pos := range positionDepths {
// 				if readStart <= uint64(pos) && readEnd >= uint64(pos) {
// 					positionDepths[pos]++
// 				}
// 			}
// 		}

// 		// if start <= 35120 && end >= 35120 && chromosome == "1" {
// 		// 	if !(rec.Flags&sam.Duplicate == sam.Duplicate) {
// 		// 		depth++
// 		// 		// fmt.Println(depth)
// 		// 		// fmt.Println(rec)
// 		// 		// fmt.Println(rec.Pos, rec.End())
// 		// 		// fmt.Println(rec.Flags)
// 		// 		// fmt.Println()
// 		// 	}
// 		// }

// 		if err != nil {
// 			log.Fatalf("error reading bam: %v", err)
// 		}
// 		// if rec.Flags&reqFlag == reqFlag && rec.Flags&excFlag == 0 {
// 		// 	n++
// 		// }
// 	}
// 	spinner.Stop()
// 	b.Close()
// }
