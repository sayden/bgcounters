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

func TestImageDraw(t *testing.T) {
	t.Run("image_draw_01", func(t *testing.T) {
		testCanvas := gg.NewContext(100, 100)

		testCanvas.Push()
		testCanvas.SetColor(color.Black)
		testCanvas.DrawRectangle(0, 0, float64(100), float64(100))
		testCanvas.Fill()
		testCanvas.Pop()

		testImage := Image{
			Settings: Settings{
				Width:        100,
				Height:       100,
				ImageScaling: IMAGE_SCALING_FIT_NONE,

				// BackgroundImage is only used for the counter and card background, but not for images
				// BackgroundImage: "",
			},
			Path: "assets/old_paper.jpg",
		}

		// Test a background image with no scaling. Settings.BackgroundImage cannot be used
		// because it is only used for the counter background
		err := testImage.Draw(testCanvas, 0, testImage.Settings)
		assert.NoError(t, err)

		// Test image scaling
		testImage.Path = "assets/binoculars.png"
		testImage.ImageScaling = IMAGE_SCALING_FIT_WIDTH
		testImage.ShadowDistance = 3
		testImage.Scale = 0.5
		testImage.Margins = 5
		err = testImage.Draw(testCanvas, 11, testImage.Settings)
		assert.NoError(t, err)

		// Test header sticked to the top
		testImage.Path = "assets/stripe.png"
		testImage.ShadowDistance = 0
		testImage.Scale = 1
		testImage.Margins = 0
		err = testImage.Draw(testCanvas, 1, testImage.Settings)
		assert.NoError(t, err)

		byt := new(bytes.Buffer)
		err = testCanvas.EncodePNG(byt)
		assert.NoError(t, err)
		assert.Equal(t, 13339, byt.Len(), "The image should have 4904 bytes but has %d bytes", byt.Len())

		// Compare the buffer with the expected image
		expectedImage, err := os.Open("testdata/image_draw_01.png")
		assert.NoError(t, err, "The expected image should be present")
		defer expectedImage.Close()

		expectedBytes, err := io.ReadAll(expectedImage)
		assert.NoError(t, err)
		assert.Equal(t, 13339, len(expectedBytes), "The expected image (from the file, not the generated image) should have 4904 bytes but has %d bytes", byt.Len())

		assert.Equal(t, expectedBytes, byt.Bytes())
	})
}
