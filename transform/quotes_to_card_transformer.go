package transform

import (
	"math/rand"

	"github.com/creasty/defaults"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/sayden/counters"
)

type QuotesToCardTransformer struct {
	Quotes         []counters.Quote
	IndexForTitles int
}

// TODO It seems that 'destination' field can be omitted
func (w *QuotesToCardTransformer) ToNewCard(cc *counters.Counter, inputCardTemplate *counters.CardsTemplate) (card *counters.Card, err error) {
	//Set a random quote
	quote := w.Quotes[rand.Intn(len(w.Quotes))]

	card, err = w.getCardPrototype()
	if err != nil {
		return nil, errors.Wrap(err, "could not read card prototype")
	}
	mergo.Merge(&card.Settings, inputCardTemplate.Settings)

	mergo.Merge(&cc.Settings, card.Settings)
	card.Areas = w.getCardAreas(cc, quote)

	return card, nil
}

func (w *QuotesToCardTransformer) getCardAreas(cc *counters.Counter, q counters.Quote) []counters.Counter {
	cc.Texts = nil
	cc.Margins = 0
	cc.BackgroundColor = ""

	// Modify images in incoming counter
	for i, img := range cc.Images {
		cc.Images[i].YShift = 0
		cc.Images[i].XShift = 0

		switch img.Position {
		case 0:
			// Use a different image
			if cc.Extra.CardImage.Path != "" {
				cc.Images[i].Path = cc.Extra.CardImage.Path
			}

			// Change the scaling
			if cc.Extra.CardImage.ImageScaling != "" {
				cc.Images[i].Path = cc.Extra.CardImage.ImageScaling
			}

			// Change the scale in case it is too small or too big
			if cc.Extra.CardImage.Scale != 0 {
				cc.Images[i].Scale = cc.Extra.CardImage.Scale
			}
		}
	}

	images, texts := w.getDownAreaCounterItems(cc, q)
	return []counters.Counter{
		*cc, {
			Images: images,
			Texts:  texts,
			Frame:  true,
		},
	}
}

func (w *QuotesToCardTransformer) getDownAreaCounterItems(cc *counters.Counter, q counters.Quote) ([]counters.Image, []counters.Text) {
	return []counters.Image{{Path: "assets/old_paper.jpg"}},
		[]counters.Text{
			{
				Settings: counters.Settings{
					AvoidClipping: true,
					StrokeWidth:   5,
					Position:      3,
				},
				String: cc.GetTextInPosition(w.IndexForTitles),
			}, {
				Settings: counters.Settings{
					FontPath:    "assets/freesans.ttf",
					StrokeWidth: 0,
					FontHeight:  20,
					FontColorS:  "black",
				},
				String: q.Quote,
			}, {
				Settings: counters.Settings{
					StrokeWidth:   0,
					FontHeight:    30,
					FontColorS:    "black",
					YShift:        -70,
					AvoidClipping: true,
					Position:      9,
				},
				String: " -" + q.Origin,
			},
		}
}

func (q *QuotesToCardTransformer) getCardPrototype() (*counters.Card, error) {
	c := counters.Card{}
	if err := defaults.Set(&c); err != nil {
		return nil, errors.Wrap(err, "could not set defaults to card")
	}

	c.BorderColorS = "black"
	c.BackgroundColor = "#faebd7"
	c.FontColorS = "white"
	c.ImageScaling = "fitHeight"
	c.StrokeColorS = "black"
	c.StrokeWidth = 0

	return &c, nil
}
