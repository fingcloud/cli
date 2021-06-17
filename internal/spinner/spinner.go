package spinner

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
)

var s *spinner.Spinner

func Start(msg ...string) {
	s = spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithColor("green"))
	if len(msg) > 0 {
		s.Suffix = fmt.Sprintf(" %s", msg[0])
	}

	s.Start()
}

func Stop(msg ...string) {
	if len(msg) > 0 {
		s.FinalMSG = msg[0] + ""
	}

	s.Stop()
}
