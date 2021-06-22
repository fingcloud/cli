package ui

import (
	"os"
	"time"

	pb "github.com/schollz/progressbar/v3"
)

func NewProgress(max int64, description ...string) *pb.ProgressBar {
	desc := ""
	if len(description) > 0 {
		desc = description[0]
	}
	bar := pb.NewOptions64(
		max,
		pb.OptionSetDescription(desc),
		pb.OptionSetWriter(os.Stderr),
		pb.OptionShowBytes(true),
		pb.OptionSetWidth(30),
		pb.OptionThrottle(65*time.Millisecond),
		pb.OptionShowCount(),
		pb.OptionSpinnerType(14),
		pb.OptionClearOnFinish(),
	)

	bar.RenderBlank()
	return bar
}
