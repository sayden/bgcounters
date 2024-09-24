package counters

import (
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/image/colornames"
	"image/color"
)

func GetValidColorForString(s string, d color.Color) color.Color {
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


func SetColorOrDefault(s, def *Settings) {
	s.BorderColor = GetValidColorForString(s.BorderColorS, def.BorderColor)
	s.FontColor = GetValidColorForString(s.FontColorS, def.FontColor)
	s.BgColor = GetValidColorForString(s.BackgroundColor, def.BgColor)
}

func SetColors(s *Settings) {
	s.BorderColor = GetValidColorForString(s.BorderColorS, color.Transparent)
	s.FontColor = GetValidColorForString(s.FontColorS, color.Transparent)
	s.BgColor = GetValidColorForString(s.BackgroundColor, color.Transparent)
	s.StrokeColor = GetValidColorForString(s.StrokeColorS, nil)
}