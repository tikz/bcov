package main

import (
	"bcov/cov"
	"bcov/db"
	"bcov/ensembl"
	"bcov/utils"
	"flag"
	"fmt"
	"log"

	"github.com/biogo/hts/sam"
	"github.com/fatih/color"

	_ "github.com/go-sql-driver/mysql"
)

var (
	require = flag.Int("f", 0, "required flags")
	exclude = flag.Int("F", 0, "excluded flags")
	file    = flag.String("file", "", "input file (empty for stdin)")
	conc    = flag.Int("threads", 0, "number of threads to use (0 = auto)")
	help    = flag.Bool("help", false, "display help")
)

func init() {
	d := color.New(color.FgRed, color.Bold)
	d.Println("Bcov")
	fmt.Println("Version: ", Version, "\t\tCommit: ", CommitHash)
	fmt.Println()
	// connectDB()

	db.ConnectDB()
}

const maxFlag = int(^sam.Flags(0))

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

func main() {
	// cov.Start()
	// start, end := db.GetChromosomeRange("1")

	// chromosome := make([end - start]uint64)
	// for i := start; i < end; i++ {
	// 	chromosome[i] = 0
	// }
	// regions := db.GetRegions()
	// fmt.Println(regions)
	// time.Sleep(8 * time.Minute)

	// for _, region := range regions {
	// 	fmt.Println(region.Chromosome, region.Start)
	// }
	cov.Load("B0434109PTC_Result/B0434109PTC.bam")
}

func main22() {
	spinner := utils.NewSpinner("Ensembl")
	spinner.Start()
	spinner.Message("connecting to public MySQL database")
	ensemblDB, err := ensembl.Connect()
	if err != nil {
		log.Fatal(err)
	}

	spinner.Message("retrieving exons...")
	exons, err := ensembl.GetExons(ensemblDB)
	if err != nil {
		log.Fatal(err)
	}
	spinner.StopDuration()

	spinner = utils.NewSpinner("Load")
	spinner.Start()
	spinner.Message("loading exons")

	geneExons := make(map[string][]ensembl.Region)
	for _, exon := range exons {
		if _, ok := geneExons[exon.GeneName]; !ok {
			geneExons[exon.GeneName] = make([]ensembl.Region, 0)
		}
		geneExons[exon.GeneName] = append(geneExons[exon.GeneName], exon)
	}

	exonCount := 0
	for geneName, exons := range geneExons {
		spinner.Message(fmt.Sprintf("%s: %d exons", geneName, len(exons)))
		geneID, created := db.StoreGene(geneName, exons[0].StableID)
		if created {
			db.StoreRegions(geneID, exons)
		}
		exonCount += len(exons)
	}

	spinner.Stop(fmt.Sprintf("%d exons in %d genes", exonCount, len(geneExons)))

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
