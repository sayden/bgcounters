package input

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonEventsToEvents(t *testing.T) {
	events, err := JsonEventsToEvents("../testdata/events.json")
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(events))
		assert.Equal(t, "Test Event title", events[0].Title)
		assert.Equal(t, "../assets/binoculars.png", events[0].Image)
		assert.Equal(t, "This is a test event", events[0].Desc)
		assert.True(t, events[0].InsertQuote)
	}
}
