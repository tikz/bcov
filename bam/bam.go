package bam

import (
	"bcov/utils"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/biogo/hts/bam"
	"github.com/biogo/hts/bgzf"
	"github.com/biogo/hts/sam"
)

type Reader struct {
	Path     string
	Filename string
	SHA256   string
	reader   *bam.Reader
}

func NewReader(path string) (*Reader, error) {
	reader := &Reader{Path: path, Filename: filepath.Base(path)}
	// reader.sha256()

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

func (r *Reader) sha256() error {
	spinner := utils.NewSpinner(fmt.Sprintf("SHA256sum %s", r.Path))
	spinner.Start()
	defer spinner.StopDuration()

	var err error
	r.SHA256, err = utils.SHA256FileHash(r.Path)
	return err
}

func (r *Reader) Read() (*sam.Record, error) {
	return r.reader.Read()
}
