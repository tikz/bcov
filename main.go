package main

import (
	"bcov/bam"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/biogo/hts/sam"

	"time"

	"github.com/theckman/yacspin"
)

var Version = "development"

var (
	require = flag.Int("f", 0, "required flags")
	exclude = flag.Int("F", 0, "excluded flags")
	file    = flag.String("file", "", "input file (empty for stdin)")
	conc    = flag.Int("threads", 0, "number of threads to use (0 = auto)")
	help    = flag.Bool("help", false, "display help")
)

const maxFlag = int(^sam.Flags(0))

func main() {
	fmt.Println("Version:\t", Version)
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *require > maxFlag {
		flag.Usage()
		log.Fatal("required flags (f) out of range")
	}
	reqFlag := sam.Flags(*require)

	if *exclude > maxFlag {
		flag.Usage()
		log.Fatal("excluded flags (F) out of range")
	}
	excFlag := sam.Flags(*exclude)

	cfg := yacspin.Config{
		Frequency:       100 * time.Millisecond,
		CharSet:         yacspin.CharSets[59],
		Suffix:          " reading file " + *file,
		SuffixAutoColon: true,
		Message:         "exporting data",
		StopCharacter:   "✓",
		StopColors:      []string{"fgGreen"},
	}

	spinner, _ := yacspin.New(cfg)
	// handle the error
	spinner.Start()

	b, err := bam.NewReader(*file)
	if err != nil {
		log.Fatal(err)
	}

	var n int
	var depth int
	for {
		rec, err := b.Read()
		if err == io.EOF {
			break
		}

		chromosome, start, end := rec.Ref.Name(), rec.Pos, rec.End()

		if start%10000 == 0 {
			spinner.Message(fmt.Sprintf("chromosome %s pos %d", chromosome, start))
		}
		if start <= 35120 && end >= 35120 && chromosome == "1" {
			if !(rec.Flags&sam.Duplicate == sam.Duplicate) {
				depth++
				// fmt.Println(depth)
				// fmt.Println(rec)
				// fmt.Println(rec.Pos, rec.End())
				// fmt.Println(rec.Flags)
				// fmt.Println()
			}
		}

		if err != nil {
			log.Fatalf("error reading bam: %v", err)
		}
		if rec.Flags&reqFlag == reqFlag && rec.Flags&excFlag == 0 {
			n++
		}
	}
	spinner.Stop()
	b.Close()

	fmt.Println(n)
}
