package ui

import (
	pb "github.com/schollz/progressbar/v3"
)

func NewProgress(max int, description ...string) *pb.ProgressBar {
	desc := ""
	if len(description) > 0 {
		desc = description[0]
	}
	bar := pb.NewOptions(
		max,
		pb.OptionShowCount(),
		pb.OptionShowBytes(true),
		pb.OptionSetDescription(desc),
		pb.OptionSpinnerType(14),
		pb.OptionEnableColorCodes(true),
		pb.OptionClearOnFinish(),
		pb.OptionSetTheme(pb.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
	bar.RenderBlank()
	return bar
}
