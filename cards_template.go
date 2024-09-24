package counters

// CardsTemplate is the template sheet (usually A4) to place cards on top in grid fashion
type CardsTemplate struct {
	Settings

	Rows    int `json:"rows,omitempty" default:"8"`
	Columns int `json:"columns,omitempty" default:"5"`

	DrawGuides bool `json:"draw_guides,omitempty"`

	// TODO is this field still used? Mode can be 'tiles' or 'template' to generate an A4 sheet like of cards or a single file per card.
	Mode string `json:"mode,omitempty" default:"tiles"`

	// TODO Rename this to OutputFolder or the one in counters to OutputPath and update JSON's
	OutputPath string `json:"output_path,omitempty" default:"output_%02d"`

	Scaling float64 `json:"scaling,omitempty" default:"1.0"`

	Cards           []Card `json:"cards"`
	MaxCardsPerFile int    `json:"max_cards_per_file"`

	// Extra Extra `json:",omitempty"`
}

// Extra is a container for extra information used in different projects but that they are not common to all of them
type Extra struct {
	FactionImage      string  `json:"faction_image,omitempty"`
	FactionImageScale float64 `json:"faction_image_scale,omitempty"`
	BackImage         string  `json:"back_image,omitempty"`
}

// ApplyCardWaterfallSettings traverses the cards in the template applying the default settings to value that are
// zero-valued
func ApplyCardWaterfallSettings(t *CardsTemplate) {
	SetColors(&t.Settings)

	for cardsIndex := range t.Cards {
		Merge(&t.Cards[cardsIndex].Settings, t.Settings)

		for counterIndex := range t.Cards[cardsIndex].Areas {
			Merge(&t.Cards[cardsIndex].Areas[counterIndex].Settings, t.Cards[cardsIndex].Settings)

			for imageIndex := range t.Cards[cardsIndex].Areas[counterIndex].Images {
				Merge(
					&t.Cards[cardsIndex].Areas[counterIndex].Images[imageIndex].Settings,
					t.Cards[cardsIndex].Areas[counterIndex].Settings,
				)
			}

			for textIndex := range t.Cards[cardsIndex].Areas[counterIndex].Texts {
				Merge(
					&t.Cards[cardsIndex].Areas[counterIndex].Texts[textIndex].Settings,
					t.Cards[cardsIndex].Areas[counterIndex].Settings,
				)
			}
		}

		for imageIndex := range t.Cards[cardsIndex].Images {
			Merge(&t.Cards[cardsIndex].Images[imageIndex].Settings, t.Cards[cardsIndex].Settings)
		}

		for textIndex := range t.Cards[cardsIndex].Texts {
			Merge(&t.Cards[cardsIndex].Texts[textIndex].Settings, t.Cards[cardsIndex].Settings)
		}
	}
}
