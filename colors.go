package counters

import (
	"image/color"

	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/image/colornames"
)

func ColorFromStringOrDefault(s string, d color.Color) color.Color {
	//Pretty color name like 'black' or 'yellow'
	col, ok := colornames.Map[s]
	if ok {
		return col
	}

	//Not recognized, maybe an hexadecimal
	hex, err := colorful.Hex(s)
	if err == nil {
		return hex
	}

	//If color is not recognized, return d or black if d is nil
	if d != nil {
		return d
	}

	return colornames.Black
}

func SetColors(s *Settings) {
	s.BorderColor = ColorFromStringOrDefault(s.BorderColorS, color.Transparent)
	s.FontColor = ColorFromStringOrDefault(s.FontColorS, color.Transparent)
	s.BgColor = ColorFromStringOrDefault(*s.BackgroundColor, color.Transparent)
	s.StrokeColor = ColorFromStringOrDefault(s.StrokeColorS, color.Transparent)
}
