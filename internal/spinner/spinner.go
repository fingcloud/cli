package spinner

import (
	"time"

	"github.com/briandowns/spinner"
)

type Spinner struct {
	*spinner.Spinner
	text string
}

func New() *Spinner {
	return &Spinner{
		Spinner: spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithColor("green")),
	}
}

func (s *Spinner) Start() *Spinner {
	s.Spinner.Suffix = " " + s.text
	s.Spinner.Stop()
	s.Spinner.Start()
	s.text = ""

	return s
}

func (s *Spinner) Text(msg string) *Spinner {
	s.text = msg
	return s
}

func (s *Spinner) Stop() *Spinner {
	s.Spinner.FinalMSG = s.text
	s.Spinner.Stop()
	s.text = ""

	return s
}
