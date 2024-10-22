package counters

import (
	"image/color"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCanvas(t *testing.T) {
	width := 800
	height := 600
	settings := &Settings{
		FontPath:   "assets/freesans.ttf",
		FontHeight: 12,
		BgColor:    color.RGBA{255, 255, 255, 255},
		Width:      width,
		Height:     height,
		FontColorS: "#000000",
	}
	template := &CardsTemplate{
		Settings: *settings,
	}

	canvas, err := GetCanvas(settings, width, height, template)
	assert.NoError(t, err)
	assert.NotNil(t, canvas)
}

func TestToImage(t *testing.T) {
	byt, err := os.ReadFile("./testdata/card_template.json")
	if !assert.NoError(t, err) {
		t.Fatal(err)
	}

	template, err := ParseCardTemplate(byt)
	if !assert.NoError(t, err) {
		t.Fatal(err)
	}

	f, _ := os.Create("/tmp/test.png")
	defer f.Close()
	err = template.Cards[0].EncodeImage(f, template)
	assert.NoError(t, err)
}

func TestApplyCardScaling(t *testing.T) {
	template := &CardsTemplate{
		Scaling: 1.5,
		Settings: Settings{
			Width:       800,
			Height:      600,
			StrokeWidth: 2,
		},
		Cards: []Card{
			{
				Settings: Settings{
					Width:  400,
					Height: 300,
				},
				Areas: []Counter{
					{Settings: Settings{Width: 100, Height: 50}},
				},
				Images: []Image{
					{Settings: Settings{Width: 200, Height: 100}},
				},
				Texts: []Text{
					{Settings: Settings{Width: 300, Height: 150, StrokeWidth: 2}},
				},
			},
		},
	}

	ApplyCardScaling(template)

	assert.Equal(t, 1200, template.Settings.Width)
	assert.Equal(t, 900, template.Settings.Height)
	assert.Equal(t, 400, template.Cards[0].Settings.Width)
	assert.Equal(t, 300, template.Cards[0].Settings.Height)
	assert.Equal(t, 150, template.Cards[0].Areas[0].Width)
	assert.Equal(t, 75, template.Cards[0].Areas[0].Height)
	assert.Equal(t, float64(3), template.Settings.StrokeWidth)
	assert.Equal(t, 300, template.Cards[0].Images[0].Width)
	assert.Equal(t, 150, template.Cards[0].Images[0].Height)
	assert.Equal(t, 450, template.Cards[0].Texts[0].Settings.Width)
	assert.Equal(t, 225, template.Cards[0].Texts[0].Settings.Height)
	assert.Equal(t, float64(3), template.Cards[0].Texts[0].Settings.StrokeWidth)
}
