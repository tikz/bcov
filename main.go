package main

import (
	"flag"
)

// List of available CLI commands
var (
	testDb        = flag.Bool("test-db", false, "Test database connection")
	fetchExons    = flag.Bool("fetch-exons", false, "Fetch exons from Ensembl and store in DB")
	fetchVariants = flag.Bool("fetch-variants", false, "Fetch SNPs from ClinVar and store in DB")
	web           = flag.Bool("web", false, "Run web server")
	bam           = flag.String("bam", "", "Load a BAM file into the database")
	bams          = flag.String("bams", "", "Load multiple BAM files passing a CSV file that in each line contains <bam path>,<capture kit name>")
	kit           = flag.String("kit", "", "Capture kit name")
	deleteBam     = flag.String("delete-bam", "", "Delete BAM file from the database with SHA256")
	region        = flag.String("region", "", "Print per position depth of a given range, expressed as <chromosome>:<start>-<end>")
	help          = flag.Bool("help", false, "Display help")
)

func init() {
	splash()
}

func main() {
	flag.Parse()
	parseFlags()
}
