package utils

import (
	"fmt"
	"time"

	"github.com/theckman/yacspin"
)

type Spinner struct {
	spinner   *yacspin.Spinner
	startTime time.Time
}

func NewSpinner(suffix string) *Spinner {
	cfg := yacspin.Config{
		Frequency:       100 * time.Millisecond,
		CharSet:         yacspin.CharSets[59],
		Suffix:          " " + suffix,
		SuffixAutoColon: true,
		StopCharacter:   "✓",
		StopColors:      []string{"fgGreen"},
	}

	spinner, _ := yacspin.New(cfg)

	return &Spinner{spinner: spinner}
}

func (s *Spinner) Message(msg string) {
	s.spinner.Message(msg)
}

func (s *Spinner) Start() {
	s.startTime = time.Now()
	s.spinner.Start()
}

func (s *Spinner) Stop(msg string) {
	s.spinner.StopMessage(msg)
	s.spinner.Stop()
}

func (s *Spinner) StopDuration() {
	duration := time.Now().Sub(s.startTime)
	s.spinner.StopMessage(fmt.Sprintf("done in %.0fm %.0fs", duration.Minutes(), duration.Seconds()))
	s.spinner.Stop()
}
