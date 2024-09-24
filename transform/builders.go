package transform

import (
	"math/rand"
	"time"

	"github.com/creasty/defaults"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"github.com/sayden/counters"
)

// DecorateBuilder is used to compose setups of counters.Counterbuilder to execute sequentially
// 'first' will be executed first in the pipelines, then the output of 'first' will be passed
// to 'second'
func DecorateBuilder(first, second counters.CounterBuilder) counters.CounterBuilder {
	return &counterBuilderDecorator{
		first:  first,
		second: second,
	}
}

type counterBuilderDecorator struct {
	first  counters.CounterBuilder
	second counters.CounterBuilder
}

func (c *counterBuilderDecorator) ToCounter(counter *counters.Counter) (*counters.Counter, error) {
	newCounter, err := c.first.ToCounter(counter)
	if err != nil {
		return nil, err
	}

	return c.second.ToCounter(newCounter)
}

// StepLossBackCounterBuilder is the most basic type of back counter builder, by adding a simple red stripe in the
// background of the counter
type StepLossBackCounterBuilder struct {
	finalOutputCounters []*counters.Counter
}

func (bc *StepLossBackCounterBuilder) ToCounter(counter *counters.Counter) (*counters.Counter, error) {
	if counter.SingleStep {
		counter.Texts = nil
		counter.Images = nil
		return counter, nil
	}

	counter.Images = append(
		[]counters.Image{
			{
				Path: counters.STRIPE, Settings: counters.Settings{ImageScaling: "fitWidth"},
			},
		}, counter.Images...,
	)

	return counter, nil
}

// SimpleFowCounterBuilder applies a default set of modification to the counters like removing numeric values
type SimpleFowCounterBuilder struct{}

func (d *SimpleFowCounterBuilder) ToCounter(cc *counters.Counter) (*counters.Counter, error) {
	if cc.Extra.PublicIcon.Path == "" {
		// No public image, no fow counter
		return cc, nil
	}

	cc.Texts = nil

	// Don't copy specific images in counters to the Fow. Take only center image, shield (if any) and air units faction
	// This way you can avoid a fow counter with infantry unit but the brigade/division icon with it. Or the flamethrower
	validFowImagesInCounter := make([]counters.Image, 0)
	for _, image := range cc.Images {
		if image.Position == 0 {
			image.Path = cc.Extra.PublicIcon.Path
			image.Scale = cc.Extra.PublicIcon.Scale
			image.YShift = 0
			image.XShift = 0
			validFowImagesInCounter = append(validFowImagesInCounter, image)
			continue
		}
	}

	cc.Images = validFowImagesInCounter

	return cc, nil
}

type QuotesToCardBuilder struct {
	Quotes         []counters.Quote
	IndexForTitles int
}

// TODO It seems that 'destination' field can be omitted
func (w *QuotesToCardBuilder) ToCard(cc counters.Counter, inputCardTemplate *counters.CardsTemplate) (card *counters.Card, err error) {
	//Set a random quote
	rand.Seed(time.Now().UnixNano())
	quote := w.Quotes[rand.Intn(len(w.Quotes))]

	card, err = getCardPrototype()
	if err != nil {
		return nil, errors.Wrap(err, "could not read card prototype")
	}
	mergo.Merge(&card.Settings, inputCardTemplate.Settings)

	mergo.Merge(&cc.Settings, card.Settings)
	card.Areas = w.getCardAreas(cc, quote)

	return card, nil
}

func (w *QuotesToCardBuilder) getCardAreas(cc counters.Counter, q counters.Quote) []counters.Counter {
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
		cc, {
			Images: images,
			Texts:  texts,
			Frame:  true,
		},
	}
}

func (w *QuotesToCardBuilder) getDownAreaCounterItems(cc counters.Counter, q counters.Quote) ([]counters.Image, []counters.Text) {
	return []counters.Image{
			{
				Path: "assets/old_paper.jpg",
			},
		},
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
					FontPath:    "/usr/share/fonts/TTF/VeraMoIt.ttf",
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

func getCardPrototype() (*counters.Card, error) {
	c := &counters.Card{}
	if err := defaults.Set(&c); err != nil {
		return nil, errors.Wrap(err, "could not set defaults to card")
	}

	c.BorderColorS = "black"
	c.BackgroundColor = "#faebd7"
	c.FontColorS = "white"
	c.ImageScaling = "fitHeight"
	c.StrokeColorS = "black"
	c.StrokeWidth = 0

	return c, nil
}
