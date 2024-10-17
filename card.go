package counters

import (
	"github.com/fogleman/gg"
	"github.com/pkg/errors"
)

type Card struct {
	Settings
	Areas  []Counter `json:"areas"`
	Texts  []Text    `json:"texts"`
	Images []Image   `json:"images"`
}

// GetCanvas returns a Canvas with attributes (like background color or size)
// taken from `settings`
func GetCanvas(settings *Settings, width, height int, t *CardsTemplate) (*gg.Context, error) {
	dc := gg.NewContext(width, height)
	err := dc.LoadFontFace(settings.FontPath, settings.FontHeight)
	if err != nil {
		return nil, errors.Wrap(err, "could not load font face")
	}

	if settings.BgColor != nil {
		dc.Push()
		dc.SetColor(settings.BgColor)
		dc.DrawRectangle(0, 0, float64(settings.Width), float64(settings.Height))
		dc.Fill()
		dc.Pop()
	}

	if settings.FontColorS != "" {
		GetValidColorForString(settings.FontColorS, t.BgColor)
	}

	return dc, nil
}

func applyCardScaling(t *CardsTemplate) {
	applySettingsScaling(&t.Settings, t.Scaling)

	for i := range t.Cards {
		c := t.Cards[i]
		applySettingsScaling(&c.Settings, t.Scaling)

		for j := range c.Areas {
			applyCounterScaling(&c.Areas[j], t.Scaling)
		}

		for j := range c.Images {
			applyImageScaling(&c.Images[j], t.Scaling)
		}

		for j := range c.Texts {
			applySettingsScaling(&c.Texts[j].Settings, t.Scaling)
		}
	}
}
