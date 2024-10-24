package counters

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/fogleman/gg"
	"github.com/stretchr/testify/assert"
)

func TestImageContent(t *testing.T, expectedImagePath string, expectedImageLength int, canvas *gg.Context) {
	byt := new(bytes.Buffer)
	err := canvas.EncodePNG(byt)
	assert.NoError(t, err)
	if assert.Equal(t, expectedImageLength, byt.Len(), "The image should have %d bytes but has %d bytes",
		expectedImageLength, byt.Len()) {

		// Compare the buffer with the expected image
		expectedImage, err := os.Open(expectedImagePath)
		if assert.NoError(t, err, "The expected image should be present") {
			defer expectedImage.Close()

			expectedBytes, err := io.ReadAll(expectedImage)
			if assert.NoError(t, err) {
				if assert.Equal(t, expectedImageLength, len(expectedBytes), "The expected image (from the file, not the "+
					"generated image) should have %d bytes but has %d bytes", expectedImageLength, byt.Len()) {

					assert.ElementsMatch(t, expectedBytes, byt.Bytes())
				}
			}
		}
	}

}

func stringP(s string) *string {
	return &s
}
