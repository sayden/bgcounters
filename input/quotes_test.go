package input

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadQuotesFromFile(t *testing.T) {
	quotes, err := ReadQuotesFromFile("../testdata/quotes.json")
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(quotes))
		assert.Equal(t, "Test Quote", quotes[0].Quote)
		assert.Equal(t, "Test Origin", quotes[0].Origin)
	}
}
