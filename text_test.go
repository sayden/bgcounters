package counters

import (
	"image/color"
	"testing"

	"github.com/fogleman/gg"
	"github.com/stretchr/testify/assert"
)

func TestTextDraw(t *testing.T) {
	t.Run("text_draw_01", func(t *testing.T) {
		sideSize := 300
		testText := Text{
			Settings: Settings{
				FontPath:   "assets/freesans.ttf",
				FontColor:  color.Black,
				FontHeight: 30,
				Width:      sideSize,
				Height:     sideSize,
				Margins:    3,
			},
			String: "11-Hello",
		}
		testCanvas := gg.NewContext(sideSize, sideSize)

		testCanvas.Push()
		testCanvas.SetColor(color.White)
		testCanvas.DrawRectangle(0, 0, float64(sideSize), float64(sideSize))
		testCanvas.Fill()
		testCanvas.Pop()

		err := testText.Draw(testCanvas, 11, testText.Settings)
		assert.NoError(t, err)

		// Test underline and different color
		testText.Underline = true
		testText.Settings.FontColor = color.RGBA{0, 0, 255, 255}
		testText.String = "3-World"
		err = testText.Draw(testCanvas, 3, testText.Settings)
		assert.NoError(t, err)

		// Text is too long for the assigned space
		testText.Settings.FontHeight = 25
		testText.String = "7-Without Clipping"
		err = testText.Draw(testCanvas, 7, testText.Settings)
		assert.NoError(t, err)

		// Ignore that text is too long
		testText.FontColor = color.White
		testText.Underline = false
		testText.String = "16-With clipping"
		testText.Settings.AvoidClipping = true
		testText.Settings.StrokeWidth = 3
		testText.StrokeColor = color.Black
		err = testText.Draw(testCanvas, 16, testText.Settings)
		assert.NoError(t, err)

		// Shadow
		testText.FontColor = color.Black
		testText.Settings.AvoidClipping = false
		testText.Settings.StrokeWidth = 0
		testText.Settings.FontColor = color.Black
		testText.ShadowDistance = 2
		testText.ShadowSigma = 2
		testText.String = "15-Sh"
		err = testText.Draw(testCanvas, 15, testText.Settings)
		assert.NoError(t, err)

		// Background color for text is red and shadows
		testText.Settings.AvoidClipping = true
		testText.TextBackgroundColor = "red"
		testText.FontHeight = 25
		testText.Settings.FontColor = color.RGBA{255, 255, 255, 255}
		testText.String = "14-Red bg"
		err = testText.Draw(testCanvas, 14, testText.Settings)

		testImageContent(t, "testdata/text_draw_01.png", 21089, testCanvas)
	})
}

func TestDrawTextsOnCanvas(t *testing.T) {
	t.Run("text_draw_01", func(t *testing.T) {
		sideSize := 300
		texts := Texts{
			{
				Settings: Settings{
					Position:   11,
					FontPath:   "assets/freesans.ttf",
					FontColorS: "black",
					FontHeight: 30,
					Width:      sideSize,
					Height:     sideSize,
					Margins:    3,
				},
				String: "11-Hello",
			},
			{
				Settings: Settings{
					Position:   3,
					FontPath:   "assets/freesans.ttf",
					FontHeight: 30,
					Width:      sideSize,
					Height:     sideSize,
					Margins:    3,
					FontColorS: "blue",
				},
				Underline: true,
				String:    "3-World",
			},
			{
				Settings: Settings{
					Position:   7,
					FontPath:   "assets/freesans.ttf",
					FontHeight: 25,
					Width:      sideSize,
					Height:     sideSize,
					Margins:    3,
					FontColorS: "blue",
				},
				Underline: true,
				String:    "7-Without Clipping",
			},
			{
				Settings: Settings{
					Position:      16,
					FontPath:      "assets/freesans.ttf",
					FontHeight:    25,
					Width:         sideSize,
					Height:        sideSize,
					Margins:       3,
					FontColorS:    "white",
					AvoidClipping: true,
					StrokeWidth:   3,
					StrokeColorS:  "black",
				},
				Underline: false,
				String:    "16-With clipping",
			},
			{
				Settings: Settings{
					Position:       15,
					FontPath:       "assets/freesans.ttf",
					FontHeight:     25,
					Width:          sideSize,
					Height:         sideSize,
					Margins:        3,
					FontColorS:     "black",
					StrokeColorS:   "black",
					ShadowDistance: 2,
					ShadowSigma:    2,
				},
				Underline: false,
				String:    "15-Sh",
			},
			{
				Settings: Settings{
					Position:       14,
					FontPath:       "assets/freesans.ttf",
					FontHeight:     25,
					Width:          sideSize,
					Height:         sideSize,
					Margins:        3,
					FontColorS:     "white",
					AvoidClipping:  true,
					StrokeColorS:   "black",
					ShadowDistance: 2,
					ShadowSigma:    2,
				},
				TextBackgroundColor: "red",
				Underline:           false,
				String:              "14-Red bg",
			},
		}

		testCanvas := gg.NewContext(sideSize, sideSize)

		testCanvas.Push()
		testCanvas.SetColor(color.White)
		testCanvas.DrawRectangle(0, 0, float64(sideSize), float64(sideSize))
		testCanvas.Fill()
		testCanvas.Pop()

		err := texts.DrawTextsOnCanvas(Settings{}, testCanvas, sideSize, sideSize)
		if assert.NoError(t, err) {
			testImageContent(t, "./testdata/text_draw_01.png", 21089, testCanvas)
		}
	})
}
