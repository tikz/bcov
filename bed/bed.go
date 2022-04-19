package bed

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
)

type Region struct {
	Chromosome string
	Start      uint64
	End        uint64
	Name       string
}

type Reader struct {
	file    *os.File
	scanner *bufio.Scanner
}

func NewReader(file string) (*Reader, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	return &Reader{file: f, scanner: bufio.NewScanner(f)}, nil
}

func (r *Reader) Read() (region Region, err error) {
	if r.file == nil {
		return Region{}, io.EOF
	}

	if r.scanner.Scan() {
		chromosome, start, end, name, errL := parseLine(r.scanner.Text())
		if errL != nil {
			return Region{}, errL
		}

		chromosome = strings.ReplaceAll(chromosome, "chr", "")
		region = Region{
			Chromosome: chromosome,
			Start:      start,
			End:        end,
			Name:       name,
		}
		return
	}

	if err := r.scanner.Err(); err != nil {
		return Region{}, err
	}

	return Region{}, io.EOF
}

func parseLine(line string) (chromosome string, start uint64, end uint64, name string, err error) {
	fields := strings.Fields(line)
	chromosome = fields[0]

	start, err = strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		return
	}

	end, err = strconv.ParseUint(fields[2], 10, 64)
	if err != nil {
		return
	}

	name = strings.Join(fields[3:], " ")
	return
}
