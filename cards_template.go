package counters

import (
	"encoding/json"

	"github.com/fogleman/gg"
	"github.com/pkg/errors"
)

// CardsTemplate is the template sheet (usually A4) to place cards on top in grid fashion
type CardsTemplate struct {
	Settings

	Rows    int `json:"rows,omitempty" default:"8"`
	Columns int `json:"columns,omitempty" default:"5"`

	DrawGuides bool `json:"draw_guides,omitempty"`

	// TODO is this field still used? Mode can be 'tiles' or 'template' to generate an A4 sheet
	// like of cards or a single file per card.
	Mode string `json:"mode,omitempty" default:"tiles"`

	// TODO Rename this to OutputFolder or the one in counters to OutputPath and update JSON's
	OutputPath string `json:"output_path,omitempty" default:"output_%02d"`

	Scaling float64 `json:"scaling,omitempty" default:"1.0"`

	Cards           []Card `json:"cards"`
	MaxCardsPerFile int    `json:"max_cards_per_file,omitempty"`

	Extra CardsExtra `json:",omitempty"`
}

// CardsExtra is a container for extra information used in different projects but that they are not
// common to all of them
type CardsExtra struct {
	FactionImage      string  `json:"faction_image,omitempty"`
	FactionImageScale float64 `json:"faction_image_scale,omitempty"`
	BackImage         string  `json:"back_image,omitempty"`
}

func ParseCardTemplate(byt []byte) (*CardsTemplate, error) {
	err := ValidateSchemaBytes[CardsTemplate](byt)
	if err != nil {
		return nil, errors.Wrap(err, "JSON file is not valid")
	}

	t := CardsTemplate{}
	if err = json.Unmarshal(byt, &t); err != nil {
		return nil, err
	}

	if t.Scaling > 0 {
		t.Settings.ApplySettingsScaling(t.Scaling)
	}

	err = t.ApplyCardWaterfallSettings()
	if err != nil {
		return nil, errors.Wrap(err, "could not apply card waterfall settings")
	}

	return &t, nil
}

// ApplyCardWaterfallSettings traverses the cards in the template applying the default settings to
// value that are zero-valued
func (t *CardsTemplate) ApplyCardWaterfallSettings() error {
	SetColors(&t.Settings)

	for cardIdx := range t.Cards {
		card := &t.Cards[cardIdx]
		if t.Scaling > 0 {
			card.ApplySettingsScaling(t.Scaling)
		}
		err := Mergev2(&card.Settings, &t.Settings)
		if err != nil {
			return err
		}

		for areaIdx := range card.Areas {
			area := &card.Areas[areaIdx]
			if t.Scaling > 0 {
				area.ApplySettingsScaling(t.Scaling)
			}
			err := Mergev2(&area.Settings, &card.Settings)
			if err != nil {
				return err
			}

			for imageIdx := range area.Images {
				image := &area.Images[imageIdx]
				if t.Scaling > 0 {
					image.ApplySettingsScaling(t.Scaling)
				}
				err := Mergev2(&image.Settings, &area.Settings)
				if err != nil {
					return err
				}
			}

			for textIdx := range area.Texts {
				text := &area.Texts[textIdx]
				if t.Scaling > 0 {
					text.ApplySettingsScaling(t.Scaling)
				}
				err := Mergev2(&text.Settings, &area.Settings)
				if err != nil {
					return err
				}
			}
		}

		for imageIdx := range card.Images {
			image := &card.Images[imageIdx]
			if t.Scaling > 0 {
				image.ApplySettingsScaling(t.Scaling)
			}
			err := Mergev2(&image.Settings, &card.Settings)
			if err != nil {
				return err
			}
		}

		for textIdx := range card.Texts {
			text := &card.Texts[textIdx]
			if t.Scaling > 0 {
				text.ApplySettingsScaling(t.Scaling)
			}
			err := Mergev2(&text.Settings, &card.Settings)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *CardsTemplate) SheetCanvas() (*gg.Context, error) {
	width := t.Columns * t.Width
	height := t.Rows * t.Height
	return t.Canvas(&t.Settings, width, height)
}

// Canvas returns a Canvas with attributes (like background color or size)
// taken from `settings`
func (t *CardsTemplate) Canvas(settings *Settings, width, height int) (*gg.Context, error) {
	dc := gg.NewContext(width, height)
	err := dc.LoadFontFace(settings.FontPath, settings.FontHeight)
	if err != nil {
		return nil, errors.Wrap(err, "could not load font face")
	}

	if settings.BgColor != nil {
		dc.Push()
		dc.SetColor(settings.BgColor)
		dc.DrawRectangle(0, 0, float64(width), float64(height))
		dc.Fill()
		dc.Pop()
	}

	if settings.FontColorS != "" {
		ColorFromStringOrDefault(settings.FontColorS, t.BgColor)
	}

	return dc, nil
}
