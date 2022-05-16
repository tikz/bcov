package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/biogo/hts/sam"
	"github.com/fatih/color"
)

var (
	testDb     = flag.Bool("test-db", false, "Test database connection")
	fetchExons = flag.Bool("fetch-exons", false, "Fetch exons from Ensembl and store in DB")
	web        = flag.Bool("web", false, "Run web server")
	bam        = flag.String("bam", "", "Load a BAM file into the database")
	bams       = flag.String("bams", "", "Load multiple BAM files passing a CSV file that in each line contains <bam path>,<capture kit name>")
	kit        = flag.String("kit", "", "Capture kit name")
	deleteBam  = flag.String("delete-bam", "", "Delete BAM file from the database with SHA256")
	region     = flag.String("region", "", "Print per position depth of a given range, expressed as <chromosome>:<start>-<end>")
	help       = flag.Bool("help", false, "Display help")
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

	if *fetchExons {
		cliFetchExons()
		os.Exit(0)
	}

	if *region != "" {
		cliExon()
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
		runWebServer()
		os.Exit(0)
	}

	fmt.Println("Usage:")
	flag.PrintDefaults()
}
