package counters

import "fmt"

// Counter is POGO-like holder for data needed for other parts to fill and draw
// a counter in a container
type Counter struct {
	Settings

	SingleStep bool `json:"single_step"`
	Frame      bool `json:"frame"`

	Images []Image `json:"images"`
	Texts  []Text  `json:"texts,omitempty"`
	Extra  *Extra  `json:"extra,omitempty"`

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

func (c *Counter) GetCounterFilename(position int, filenumber int, suffix string) string {
	name := c.GetTextInPosition(position)

	if filenumber >= 0 {
		if position >= 0 && position <= 14 {
			if c.Extra != nil {
				return fmt.Sprintf("%s_%s_%03d%s.png",
					c.Extra.Side, name, filenumber, suffix)
			}
			if suffix != "" {
				return fmt.Sprintf("%03d_%s_%s.png", filenumber, name, suffix)
			}

			return fmt.Sprintf("%03d_%s.png", filenumber, name)
		}
		return fmt.Sprintf("%03d%s.png", filenumber, suffix)
	}

	return fmt.Sprintf("%s_%s%s.png", c.Extra.Side, name, suffix)
}

type CounterPrototype struct {
	Counter

	ImagesPrototypes []ImagePrototype `json:"image_prototypes,omitempty"`
	TextsPrototypes  []TextPrototype  `json:"text_prototypes,omitempty"`
}

func applyCounterScaling(c *Counter, scaling float64) {
	for i := range c.Images {
		applyImageScaling(&c.Images[i], scaling)
	}

	for i := range c.Texts {
		applySettingsScaling(&c.Texts[i].Settings, scaling)
	}
}
