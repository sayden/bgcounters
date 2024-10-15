package counters

import (
	"bytes"
	"image/color"
	"io"
	"os"
	"testing"

	"github.com/fogleman/gg"
	"github.com/stretchr/testify/assert"
)

func TestTextDraw(t *testing.T) {
	t.Run("text_draw_01", func(t *testing.T) {
		sideSize := 300
		testCanvas := gg.NewContext(sideSize, sideSize)

		testCanvas.Push()
		testCanvas.SetColor(color.White)
		testCanvas.DrawRectangle(0, 0, float64(sideSize), float64(sideSize))
		testCanvas.Fill()
		testCanvas.Pop()

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

		byt := new(bytes.Buffer)
		err = testCanvas.EncodePNG(byt)
		assert.NoError(t, err)
		assert.Equal(t, 21089, byt.Len(), "The image should have 21089 bytes but has %d bytes", byt.Len())

		// Compare the buffer with the expected image
		expectedImage, err := os.Open("testdata/text_draw_01.png")
		assert.NoError(t, err, "The expected image should be present")
		defer expectedImage.Close()

		expectedBytes, err := io.ReadAll(expectedImage)
		assert.NoError(t, err)
		assert.Equal(t, 21089, len(expectedBytes), "The expected image (from the file, not the "+
			"generated image) should have 21089 bytes but has %d bytes", byt.Len())

		assert.Equal(t, expectedBytes, byt.Bytes())
	})
}
