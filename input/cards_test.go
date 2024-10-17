package input

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadJSONCardsFile(t *testing.T) {
	// read the json file with testing data
	filepath := "../testdata/cards.json"
	cardTemplate, err := ReadJSONCardsFile(filepath)
	assert.NoError(t, err, "could not read the cards file")
	assert.Equal(t, 1, len(cardTemplate.Cards))
	assert.Equal(t, "Card Title here", cardTemplate.Cards[0].Texts[0].String)
	assert.Equal(t, "../assets/binoculars.png", cardTemplate.Cards[0].Areas[0].Images[0].Path)
	assert.Equal(t, 742, cardTemplate.Cards[0].Settings.Width)
	assert.Equal(t, 1200, cardTemplate.Cards[0].Settings.Height)
}
