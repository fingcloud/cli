package ui

import "github.com/logrusorgru/aurora/v3"

var (
	color = aurora.NewAurora(true)
)

func Bold(v interface{}) aurora.Value   { return color.Bold(v) }
func Green(v interface{}) aurora.Value  { return color.Green(v) }
func Yellow(v interface{}) aurora.Value { return color.Yellow(v) }
func Red(v interface{}) aurora.Value    { return color.Red(v) }
func Blue(v interface{}) aurora.Value   { return color.BrightBlue(v) }
func Gray(v interface{}) aurora.Value   { return color.Gray(14, v) }

func Heading(v string) string {
	return color.Sprintf(color.Bold("==> %s\n").Magenta(), Bold(v))
}

func Alert(v string) string {
	return color.Sprintf("%s %s", color.Red("==>"), color.Bold(v))
}

func Warning(v string) string {
	return color.Sprintf("%s %s", color.Yellow("==>"), color.Bold(v))
}

func Info(v string) string {
	return color.Sprintf("%s %s", color.Blue("==>"), color.Bold(v))
}

func KeyValue(key string, value interface{}) string {
	return color.Sprintf("%s %v", key, value)
}
