package view

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"myspoti/internal/player"
)

const (
	boxInnerWidth  = 46
	progressBarLen = 22
)

// Playback renders the now-playing card.
func Playback(p player.Playback) string {
	bars := "▍▊▌▉▎"
	title := truncateRunes(p.Track, boxInnerWidth-11)
	subtitle := truncateRunes(fmt.Sprintf("%s · %s", p.Artists, p.Album), boxInnerWidth-9)

	playIcon := "⏸"
	if p.IsPlaying {
		playIcon = "▶"
	}

	progress := formatDuration(p.Progress)
	duration := formatDuration(p.Duration)
	bar := progressBar(p.Progress, p.Duration, progressBarLen)

	line1 := fmt.Sprintf("  %s  %s", bars, title)
	line2 := fmt.Sprintf("         %s", subtitle)
	line3 := ""
	line4 := fmt.Sprintf("  %s  %s %s %s", playIcon, progress, bar, duration)

	var b strings.Builder
	writeBorder(&b, '╭', '╮', true)
	writeBoxLine(&b, line1)
	writeBoxLine(&b, line2)
	writeBoxLine(&b, line3)
	writeBoxLine(&b, line4)
	writeBorder(&b, '╰', '╯', false)
	return b.String()
}

func formatDuration(d time.Duration) string {
	total := int(d.Seconds())
	if total < 0 {
		total = 0
	}
	return fmt.Sprintf("%d:%02d", total/60, total%60)
}

func writeBorder(b *strings.Builder, left, right rune, newline bool) {
	b.WriteRune(left)
	b.WriteString(strings.Repeat("─", boxInnerWidth))
	b.WriteRune(right)
	if newline {
		b.WriteByte('\n')
	}
}

func writeBoxLine(b *strings.Builder, content string) {
	pad := boxInnerWidth - utf8.RuneCountInString(content)
	if pad < 0 {
		content = truncateRunes(content, boxInnerWidth)
		pad = 0
	}
	b.WriteRune('│')
	b.WriteString(content)
	b.WriteString(strings.Repeat(" ", pad))
	b.WriteRune('│')
	b.WriteByte('\n')
}

func progressBar(progress, duration time.Duration, width int) string {
	if width < 1 {
		return ""
	}
	pos := 0
	if duration > 0 {
		pos = int(float64(progress) / float64(duration) * float64(width-1))
	}
	if pos < 0 {
		pos = 0
	}
	if pos > width-1 {
		pos = width - 1
	}

	var b strings.Builder
	for i := range width {
		switch {
		case i == pos:
			b.WriteRune('●')
		case i < pos:
			b.WriteRune('━')
		default:
			b.WriteRune('─')
		}
	}
	return b.String()
}

func truncateRunes(s string, max int) string {
	if max <= 0 {
		return ""
	}
	if utf8.RuneCountInString(s) <= max {
		return s
	}
	if max == 1 {
		return "…"
	}
	runes := []rune(s)
	return string(runes[:max-1]) + "…"
}
