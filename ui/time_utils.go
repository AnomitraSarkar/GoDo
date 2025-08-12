package ui

import (
	"fmt"
	"time"
)

func formatTime(t time.Time, relative bool) string {
	if t.IsZero() {
		return "unknown"
	}

	if relative {
		now := time.Now()
		diff := now.Sub(t)

		switch {
		case diff < time.Minute:
			return "just now"
		case diff < time.Hour:
			mins := int(diff.Minutes())
			if mins == 1 {
				return "1m ago"
			}
			return fmt.Sprintf("%dm ago", mins)
		case diff < 24*time.Hour:
			hours := int(diff.Hours())
			if hours == 1 {
				return "1h ago"
			}
			return fmt.Sprintf("%dh ago", hours)
		case diff < 30*24*time.Hour:
			days := int(diff.Hours() / 24)
			if days == 1 {
				return "1d ago"
			}
			return fmt.Sprintf("%dd ago", days)
		case diff < 365*24*time.Hour:
			months := int(diff.Hours() / (24 * 30))
			if months == 1 {
				return "1M ago"
			}
			return fmt.Sprintf("%dM ago", months)
		default:
			years := int(diff.Hours() / (24 * 365))
			if years == 1 {
				return "1y ago"
			}
			return fmt.Sprintf("%dy ago", years)
		}
	}
	return t.Format("02-01-2006 15:04:05")
}