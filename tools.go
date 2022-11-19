package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

func byteCountIEC(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%dB", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	num := float64(b) / float64(div)
	format := fmt.Sprintf("%.1f", num)
	if format[len(format)-2:] == ".0" {
		format = format[:len(format)-2]
	}
	return fmt.Sprintf("%s%cB", format, "KMGTPE"[exp])
}

func firstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func formatUptime(uptime time.Duration) string {
	days := math.Floor(uptime.Hours() / 24)
	uptime -= time.Hour * 24 * time.Duration(days)
	var dayStr string
	if days > 0 {
		dayStr = fmt.Sprintf("%ddays ", int(days))
	}
	return fmt.Sprintf("%s%s", dayStr, uptime)
}

func markdownEscape(text string) string {
	replacer := strings.NewReplacer(`_`, `\_`, `*`, `\*`, `[`, `\[`, "`", "\\`")
	return replacer.Replace(text)
}

func in(s string, list []string) bool {
	for _, item := range list {
		if item == s {
			return true
		}
	}
	return false
}
