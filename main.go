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
	bam          = flag.String("bam", "", "Load a BAM file into the database")
	deleteBam    = flag.String("delete-bam", "", "Delete BAM file from the database with SHA256")
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
