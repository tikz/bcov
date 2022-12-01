package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/tikz/bcov/api"
	"github.com/tikz/bcov/cov"
	"github.com/tikz/bcov/db"
	"github.com/tikz/bcov/ensembl"
	"github.com/tikz/bcov/utils"

	"github.com/tikz/bio/clinvar"
)

// splash prints the greeting and current version info.
func splash() {
	d := color.New(color.FgRed, color.Bold)
	d.Println("Bcov")
	fmt.Println("Version: ", Version, "\t\tCommit: ", CommitHash)
	fmt.Println()
}

// parseFlags starts the required action by the user, or prints the usage help.
func parseFlags() {
	if *help {
		fmt.Println("Usage:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *testDb {
		cliTestDB()
		os.Exit(0)
	}

	if *fetchExons {
		cliFetchExons()
		os.Exit(0)
	}

	if *fetchVariants {
		cliFetchVariants()
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

	if *bams != "" {
		cliLoadBAMs()
		os.Exit(0)
	}

	if *deleteBam != "" {
		cliDeleteBam()
		os.Exit(0)
	}

	if *web {
		api.RunServer()
		os.Exit(0)
	}

	fmt.Println("Usage:")
	flag.PrintDefaults()
}

// cliTestDB handles printing object counts from the database when called from the CLI
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
	db.DB.Model(&db.Exon{}).Count(&count)
	fmt.Println("Exons:", count)
	db.DB.Model(&db.ExonDepthCoverage{}).Count(&count)
	fmt.Println("Depth coverages:", count)
}

// cliRegion handles showing per depth position of a given chromosome and start end position of a BAM file when called from the CLI
// Used for debug purposes and quickly checking a BAM file content.
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

	exon := cov.ExonDepth(*bam, chromosome, start, end)
	fmt.Println("CHROMOSOME:POSITION \t DEPTH")
	for i := start; i < end; i++ {
		fmt.Printf("%s:%d \t\t\t %d\n", chromosome, i, exon.PositionDepth[cov.Position(i)])
	}

}

// cliFetchExons handles the load of exons from Ensemble to the database when called from the CLI
func cliFetchExons() {
	db.ConnectDB()

	spinner := utils.NewSpinner("Ensembl")
	spinner.Start()
	spinner.Message("connecting to public database")
	ensemblDB, err := ensembl.Connect()
	if err != nil {
		log.Fatal(err)
	}

	spinner.Message("running Ensembl query for retrieving all exons...")
	exons, err := ensembl.GetExons(ensemblDB)
	if err != nil {
		log.Fatal(err)
	}
	spinner.StopDuration()

	spinner = utils.NewSpinner("Load")
	spinner.Start()
	spinner.Message("loading exons")

	geneExons := make(map[string][]ensembl.Exon)
	for _, exon := range exons {
		if _, ok := geneExons[exon.GeneName]; !ok {
			geneExons[exon.GeneName] = make([]ensembl.Exon, 0)
		}
		geneExons[exon.GeneName] = append(geneExons[exon.GeneName], exon)
	}

	exonCount := 0
	for geneName, exons := range geneExons {
		// Store exons in database
		geneID, created := db.StoreGene(exons[0].HGNCAccession, exons[0].GeneAccession, geneName, exons[0].GeneDescription, exons[0].TranscriptAccession)
		if created {
			db.StoreExons(geneID, exons)
			spinner.Message(fmt.Sprintf("%s: %d exons", geneName, len(exons)))
		}

		// Store synonyms in database
		synonyms, err := ensembl.GetSynonyms(ensemblDB, geneName)
		if err != nil {
			log.Fatalf("Error getting synonyms for gene %s: %s", geneName, err)
		}
		created = db.StoreSynonyms(geneID, synonyms)

		if created {
			spinner.Message(fmt.Sprintf("%s: with %d gene synonyms", geneName, len(synonyms)))
		}

		exonCount += len(exons)
	}

	spinner.Stop(fmt.Sprintf("%d exons in %d genes", exonCount, len(geneExons)))
}

// cliFetchVariants handles the loading of variants from ClinVar to the database when called from the CLI
func cliFetchVariants() {
	db.ConnectDB()

	spinner := utils.NewSpinner("ClinVar")
	spinner.Start()
	spinner.Message("downloading variant_summary.txt")
	cv, err := clinvar.NewClinVar("/tmp")
	if err != nil {
		log.Fatal(err)
	}

	spinner.Message("loading variants")

	// Assign SNPs to exon IDs
	// sort first for performance reasons
	exons := db.GetExons()
	var snps []clinvar.Allele

	for _, alleles := range cv.SNPs {
		for _, allele := range alleles {
			snps = append(snps, allele)
		}
	}

	sort.SliceStable(snps, func(i, j int) bool {
		ci := utils.ChromosomeIndex(snps[i].Chromosome)
		cj := utils.ChromosomeIndex(snps[j].Chromosome)

		if ci < cj {
			return true
		}
		if ci > cj {
			return false
		}
		return snps[i].Start < snps[j].Start
	})

	count := 0
	i, j := 0, 0

	for i < len(exons) && j < len(snps) {
		exon := exons[i]
		snp := snps[j]

		if exon.Chromosome == snp.Chromosome && snp.Start >= exon.Start && snp.End <= exon.End {
			spinner.Message(fmt.Sprintf("(%.2f%%) variant: rs%s", 100*float64(count)/float64(len(snps)), snp.VariantID))
			db.StoreVariant(snp, exon.ID)
			count++
		}

		if utils.ChromosomeIndex(exons[i].Chromosome) < utils.ChromosomeIndex(snps[j].Chromosome) {
			i++
		}

		if utils.ChromosomeIndex(exons[i].Chromosome) > utils.ChromosomeIndex(snps[j].Chromosome) {
			j++
		}

		if snps[j].Start > exons[i].End {
			i++
		} else {
			j++
		}
	}

	spinner.Stop(fmt.Sprintf("%d variants", count))
}

// checkExonsDB handles showing a count of loaded exones when called from the CLI
func checkExonsDB() {
	var count int64
	db.DB.Model(&db.Exon{}).Count(&count)

	if count == 0 {
		fmt.Println("No exons found in database. You may want to run -fetch-exons first.")
		os.Exit(1)
	}
}

// cliLoadBAM handles the load of a single BAM file by path when called from the CLI
func cliLoadBAM() {
	db.ConnectDB()
	checkExonsDB()

	if *kit == "" {
		fmt.Println(`Capture kit name not specified for this bam, add -kit "name"`)
		os.Exit(1)
	}

	cov.Load(*bam, *kit)
}

// cliLoadBAMs handles the load of multiple BAM files from a CSV file when called from the CLI
func cliLoadBAMs() {
	db.ConnectDB()
	checkExonsDB()

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

// cliDeleteBam handles the deleting a BAM file from the database by SHA256sum when called from the CLI
func cliDeleteBam() {
	db.ConnectDB()

	var count int64
	db.DB.Model(&db.Exon{}).Count(&count)

	var bamFile db.BAMFile
	result := db.DB.First(&bamFile, "sha256_sum = ?", *deleteBam)
	if result.RowsAffected == 0 {
		fmt.Printf("No BAM files found in database with hash %s\n", *deleteBam)
		os.Exit(1)
	}

	db.DB.Delete(&bamFile)
	fmt.Println("Successfully deleted BAM file and associated records from the database")
}
