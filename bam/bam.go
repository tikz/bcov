package bam

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/biogo/hts/bam"
	"github.com/biogo/hts/bgzf"
)

func NewReader(path string) (*bam.Reader, error) {
	var r io.Reader
	if path == "" {
		r = os.Stdin
	} else {
		f, err := os.Open(path)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("could not open file %q:", err))
		}
		ok, err := bgzf.HasEOF(f)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("could not open file %q:", err))
		}
		if !ok {
			return nil, errors.New(fmt.Sprintf("file %q has no bgzf magic block: may be truncated", path))
		}
		r = f
	}

	b, err := bam.NewReader(r, 0)
	if err != nil {
		return nil, errors.New(fmt.Sprint("could not read bam:", err))
	}

	b.Omit(bam.AllVariableLengthData)

	return b, nil
}
