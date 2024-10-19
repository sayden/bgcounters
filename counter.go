package counters

import (
	"fmt"
	"io"
	"strings"

	"github.com/fogleman/gg"
	"github.com/pkg/errors"
	"github.com/thehivecorporation/log"
)

// Counter is POGO-like holder for data needed for other parts to fill and draw
// a counter in a container
type Counter struct {
	Settings

	SingleStep bool `json:"single_step,omitempty"`
	Frame      bool `json:"frame,omitempty"`

	Images Images `json:"images,omitempty"`
	Texts  Texts  `json:"texts,omitempty"`
	Extra  *Extra `json:"extra,omitempty"`

	// Generate the following counter with 'back' suffix in its filename
	Back *Counter `json:",omitempty"`
}

// TODO This Extra contains data from all projects
type Extra struct {
	PublicIcon         *imageExtraData `json:"public_icon,omitempty"`
	CardImage          *imageExtraData `json:"card_image,omitempty"`
	SkipCardGeneration bool            `json:"skip_card_generation,omitempty"`
	Title              string          `json:"title,omitempty"`
	Cost               int             `json:"cost,omitempty"`
	Side               string          `json:"side,omitempty"`
	TitlePosition      *int            `json:"title_position,omitempty"`
}

type imageExtraData struct {
	// the path to find the image file
	Path string `json:"path,omitempty"`

	// a percentage of original image's size
	Scale float64 `json:"scale,omitempty"`

	// none, fitHeight, fitWidth or wrap
	ImageScaling string `json:"image_scaling,omitempty"`
}

func (c *Counter) GetTextInPosition(i int) string {
	for _, text := range c.Texts {
		if text.Position == i {
			return text.String
		}
	}

	return ""
}

// filenumber: CounterTemplate.PositionNumberForFilename. So it will always be fixed number
// position: The position of the text in the counter (0-16)
// suffix: A suffix on the file. Constant
func (c *Counter) GetCounterFilename(position int, suffix string, filenumber int, filenamesInUse map[string]bool) string {
	var b strings.Builder
	name := c.GetTextInPosition(position)

	if c.Extra != nil {
		if c.Extra.TitlePosition != nil && *c.Extra.TitlePosition != position {
			name = c.GetTextInPosition(*c.Extra.TitlePosition)
		}
		if name != "" {
			b.WriteString(name + " ")
		}
		// This way, the positional based name will always be the first part of the filename
		// while the manual title will come later. This is useful when using prototypes so that
		// counters with the same positional name are close together in the destination folder
		// by "formation" (belonging) instead of by "use" (title)
		name = ""

		if c.Extra.Side != "" {
			b.WriteString(c.Extra.Side)
			b.WriteString(" ")
		}

		if c.Extra.Title != "" {
			b.WriteString(c.Extra.Title)
			b.WriteString(" ")
		}
	}

	if name != "" {
		b.WriteString(name + " ")
	}

	if suffix != "" {
		b.WriteString(suffix)
	}

	res := b.String()
	res = strings.TrimSpace(res)

	if filenamesInUse[res] {
		if filenumber >= 0 {
			res += fmt.Sprintf(" %03d", filenumber)
		}
	}

	filenamesInUse[res] = true

	res += ".png"

	return res
}

func (c *Counter) Canvas(withGuides bool) (*gg.Context, error) {
	canvas, err := c.canvas()
	if err != nil {
		return nil, err
	}

	// Draw background image
	if err = c.DrawBackgroundImage(canvas); err != nil {
		return nil, errors.Wrap(err, "error trying to draw background image")
	}

	// Draw images
	if err = c.Images.DrawImagesOnCanvas(&c.Settings, canvas, c.Width, c.Height); err != nil {
		return nil, errors.Wrap(err, "error trying to process image")
	}

	// Draw texts
	if err = c.Texts.DrawTextsOnCanvas(c.Settings, canvas, c.Width, c.Height); err != nil {
		return nil, errors.Wrap(err, "error trying to draw text")
	}

	// Draw guides
	if withGuides {
		guides, err := DrawGuides(c.Settings)
		if err != nil {
			return nil, err
		}
		canvas.DrawImage(*guides, 0, 0)
	}

	return canvas, nil
}

func (c *Counter) canvas() (*gg.Context, error) {
	dc := gg.NewContext(c.Width, c.Height)
	if err := dc.LoadFontFace(c.FontPath, c.FontHeight); err != nil {
		log.WithFields(log.Fields{"font": "'" + c.FontPath + "'", "height": c.FontHeight}).Error(err)
		return nil, err
	}

	dc.Push()
	dc.SetColor(c.BgColor)
	dc.DrawRectangle(0, 0, float64(c.Width), float64(c.Height))
	dc.Fill()
	dc.Pop()

	if c.FontColorS != "" && c.FontColor == nil {
		ColorFromStringOrDefault(c.FontColorS, c.FontColor)
	}

	return dc, nil
}

func ApplyCounterScaling(c *Counter, scaling float64) {
	for i := range c.Images {
		applyImageScaling(&c.Images[i], scaling)
	}

	for i := range c.Texts {
		ApplySettingsScaling(&c.Texts[i].Settings, scaling)
	}

	ApplySettingsScaling(&c.Settings, scaling)
}

func (c *Counter) EncodeCounter(w io.Writer, template *CounterTemplate) error {
	counterCanvas, err := c.Canvas(template.DrawGuides)
	if err != nil {
		return err
	}

	return counterCanvas.EncodePNG(w)
}
