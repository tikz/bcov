package bam

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/tikz/bcov/utils"

	"github.com/biogo/hts/bam"
	"github.com/biogo/hts/bgzf"
	"github.com/biogo/hts/sam"
)

type Reader struct {
	Path     string
	Filename string
	Size     int64
	reader   *bam.Reader
}

func NewReader(path string) (*Reader, error) {
	reader := &Reader{Path: path, Filename: filepath.Base(path)}

	fi, err := os.Stat(path)
	if err != nil {
		return nil, errors.New(fmt.Sprint("could not stat file:", err))
	}
	reader.Size = fi.Size()

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

	reader.reader = b
	return reader, nil
}

func (r *Reader) SHA256sum() string {
	spinner := utils.NewSpinner(fmt.Sprintf("SHA256sum %s", r.Path))
	spinner.Start()
	defer spinner.StopDuration()

	hash, err := utils.SHA256FileHash(r.Path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return hash
}

func (r *Reader) Read() (*sam.Record, error) {
	return r.reader.Read()
}
