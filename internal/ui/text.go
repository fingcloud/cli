package ui

import "github.com/logrusorgru/aurora/v3"

var (
	color = aurora.NewAurora(true)
)

func Bold(v interface{}) aurora.Value  { return color.Bold(v) }
func Green(v interface{}) aurora.Value { return color.Green(v) }
func Red(v interface{}) aurora.Value   { return color.Red(v) }
func Blue(v interface{}) aurora.Value  { return color.BrightBlue(v) }
func Gray(v interface{}) aurora.Value  { return color.Gray(14, v) }

func Heading(v string) string {
	return color.Sprintf(color.Bold("==> %s\n").Magenta(), Bold(v))
}

func Alert(v string) string {
	return color.Sprintf(color.Bold("🚨 %s\n").Red(), v)
}

func Warning(text string) string {
	return color.Sprintf(color.Bold("⚠️ %s\n").Yellow(), text)
}

func Info(text string) string {
	return color.Sprintf(Gray(Bold("💁 %s\n").String()), text)
}

func KeyValue(key string, value interface{}) string {
	return color.Sprintf("%s %v", key, value)
}
