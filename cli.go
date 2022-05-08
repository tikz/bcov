package main

import (
	"bcov/cov"
	"bcov/db"
	"bcov/ensembl"
	"bcov/utils"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func cliTestDB() {
	db.ConnectDB()

	fmt.Println("Connected to database.")

	fmt.Println("Engine:", db.DB.Config.Dialector.Name())
	fmt.Println()

	var count int64
	db.DB.Model(&db.BAMFile{}).Count(&count)
	fmt.Println("BAM files:", count)
	db.DB.Model(&db.Gene{}).Count(&count)
	fmt.Println("Genes:", count)
	db.DB.Model(&db.Region{}).Count(&count)
	fmt.Println("Regions:", count)
	db.DB.Model(&db.RegionDepthCoverage{}).Count(&count)
	fmt.Println("Depth coverages:", count)
}

func cliRegion() {
	if *bam == "" {
		fmt.Println("BAM file required.")
		os.Exit(1)
	}
	sepColon := strings.Split(*region, ":")
	sepDash := strings.Split(sepColon[1], "-")

	chromosome := sepColon[0]
	start, err := strconv.ParseUint(sepDash[0], 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	end, err := strconv.ParseUint(sepDash[1], 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	region := cov.RegionDepth(*bam, chromosome, start, end)
	fmt.Println("CHROMOSOME:POSITION \t DEPTH")
	for i := start; i < end; i++ {
		fmt.Printf("%s:%d \t\t\t %d\n", chromosome, i, region.PositionDepth[cov.Position(i)])
	}

}

func cliFetchRegions() {
	db.ConnectDB()

	spinner := utils.NewSpinner("Ensembl")
	spinner.Start()
	spinner.Message("connecting to public database")
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
		geneID, created := db.StoreGene(exons[0].GeneAccession, geneName, exons[0].GeneDescription, exons[0].StableID)
		if created {
			db.StoreRegions(geneID, exons)
		}
		exonCount += len(exons)
	}

	spinner.Stop(fmt.Sprintf("%d exons in %d genes", exonCount, len(geneExons)))
}

func checkRegionsDB() {
	var count int64
	db.DB.Model(&db.Region{}).Count(&count)

	if count == 0 {
		fmt.Println("No regions found in database. You may want to run -fetch-regions first.")
		os.Exit(1)
	}

}

func cliLoadBAM() {
	db.ConnectDB()
	checkRegionsDB()

	if *kit == "" {
		fmt.Println(`Capture kit name not specified for this bam, add -kit "name"`)
		os.Exit(1)
	}

	cov.Load(*bam, *kit)
}

func cliLoadBAMs() {
	db.ConnectDB()
	checkRegionsDB()

	f, err := os.Open(*bams)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		path, kitName := strings.TrimSpace(row[0]), strings.TrimSpace(row[1])
		cov.Load(path, kitName)
	}
}

func cliDeleteBam() {
	db.ConnectDB()

	var count int64
	db.DB.Model(&db.Region{}).Count(&count)

	var bamFile db.BAMFile
	result := db.DB.First(&bamFile, "sha256_sum = ?", *deleteBam)
	if result.RowsAffected == 0 {
		fmt.Printf("No BAM files found in database with hash %s\n", *deleteBam)
		os.Exit(1)
	}

	db.DB.Delete(&bamFile)
	fmt.Println("Successfully deleted BAM file and associated records from the database")
}
