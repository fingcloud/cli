package util

import (
	"fmt"
	"time"
)

func pluralize(num int, s string) string {
	if num == 1 {
		return fmt.Sprintf("%d %s", num, s)
	}
	return fmt.Sprintf("%d %ss", num, s)
}

func fmtDuration(amount int, unit string) string {
	return fmt.Sprintf("about %s ago", pluralize(amount, unit))
}

// FuzzyAgo returns humanized time ago
func FuzzyAgo(d time.Duration) string {
	if d < time.Minute {
		return "less than a minute ago"
	}
	if d < time.Hour {
		return fmtDuration(int(d.Minutes()), "minute")
	}
	if d < 24*time.Hour {
		return fmtDuration(int(d.Hours()), "hour")
	}
	if d < 24*30*time.Hour {
		return fmtDuration(int(d.Hours()/24), "day")
	}
	if d < 24*30*12*time.Hour {
		return fmtDuration(int(d.Hours()/24/30), "month")
	}
	return fmtDuration(int(d.Hours()/24/365), "year")
}

func FuzzyAgoAbbr(now time.Time, createdAt time.Time) string {
	ago := now.Sub(createdAt)

	if ago < time.Hour {
		return fmt.Sprintf("%d%s", int(ago.Minutes()), "m")
	}
	if ago < 24*time.Hour {
		return fmt.Sprintf("%d%s", int(ago.Hours()), "h")
	}
	if ago < 30*24*time.Hour {
		return fmt.Sprintf("%d%s", int(ago.Hours())/24, "d")
	}

	return createdAt.Format("Jan _2, 2006")
}
