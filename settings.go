package counters

import (
	"image/color"
)

// Settings are used in counters and cards. Represents data required to draw and position a bounding box into a container
type Settings struct {
	Width  int `json:"width,omitempty"`
	Height int `json:"height,omitempty"`

	Margins float64 `json:"margins,omitempty"`

	FontHeight float64 `json:"font_height,omitempty"`
	FontPath   string  `json:"font_path,omitempty" default:"assets/font-bebas.ttf"`

	FontColorS string      `json:"font_color,omitempty" default:"black"`
	FontColor  color.Color `json:"-"`

	BackgroundImage string      `json:"background_image,omitempty"`
	BackgroundColor string      `json:"background_color,omitempty" default:"white"`
	BgColor         color.Color `json:"-"`

	ShadowDistance int `json:"shadow,omitempty" default:"0"`

	Rotation float64 `json:"rotation,omitempty"`

	//CounterTemplate Card specific
	BorderWidth  float64     `json:"border_width,omitempty" default:"0"`
	BorderColorS string      `json:"border_color,omitempty"`
	BorderColor  color.Color `json:"-"`

	XShift float64 `json:"x_shift,omitempty"`
	YShift float64 `json:"y_shift,omitempty"`

	//Card specific
	Multiplier int `json:"multiplier,omitempty" default:"1"`

	//Text options specific
	StrokeWidth  float64     `json:"stroke_width,omitempty" default:"0"`
	StrokeColorS string      `json:"stroke_color,omitempty" default:"black"`
	StrokeColor  color.Color `json:"-"`
	Alignment    string      `json:"alignment,omitempty"`

	//Image options specific
	ImageScaling string `json:"image_scaling,omitempty" default:"none"`

	AvoidClipping bool `json:"avoid_clipping,omitempty"`

	Position int `json:"position,omitempty"`

	SkipBorders bool `json:"skip_borders,omitempty"`

	Skip bool `json:"skip,omitempty"`
}

func (i *Settings) GetPosition() int {
	return i.Position
}

func applySettingsScaling(s *Settings, scaling float64) {
	s.Width = int(scaling * float64(s.Width))
	s.Height = int(scaling * float64(s.Height))

	s.Margins *= scaling
	if s.Margins < 1 {
		s.Margins = 1
	}

	s.FontHeight *= scaling

	s.ShadowDistance = int(scaling * float64(s.ShadowDistance))

	s.BorderWidth *= scaling
	if s.BorderWidth < 1 {
		s.BorderWidth = 1
	}

	s.XShift *= scaling
	s.YShift *= scaling

	s.StrokeWidth *= scaling
}
