package input

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadCSVCounters(t *testing.T) {
	inputPathCsv := "../testdata/testing_cards.csv"

	ct, err := ReadCSVCounters(inputPathCsv)
	if err != nil {
		t.Error("Error reading the CSV file")
	}
	assert.Equal(t, 2, len(ct.Counters))
}

func TestReadCSVCards(t *testing.T) {
	inputPathCsv := "../testdata/testing_cards.csv"
	f, err := os.Open(inputPathCsv)
	if assert.NoError(t, err) {
		template, err := ReadJSONCardsFile("../testdata/cards.json")
		if assert.NoError(t, err, "could not read the cards file") {
			template.Cards = nil
			ct, err := ReadCSVCards(f, template)
			if err != nil {
				t.Error("Error reading the CSV file")
			}
			assert.Equal(t, 21, len(ct.Cards))
		}
	}

}
