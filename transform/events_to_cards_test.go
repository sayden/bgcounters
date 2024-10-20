package transform

import (
	"testing"

	"github.com/sayden/counters"
)

func TestEventsToCards(t *testing.T) {
	// Sample events
	events := []counters.Event{
		{Title: "Event 1", Desc: "Description 1"},
		{Title: "Event 2", Desc: "Description 2"},
	}

	// Sample images
	images := []string{"image1.jpg", "image2.jpg"}

	// Configuration for EventsToCards
	cfg := &EventsToCardsConfig{
		Events:             events,
		Images:             images,
		BackImageFile:      "back_image.jpg",
		GeneratedImageName: "output.png",
	}

	// Call the function
	cardTemplate := EventsToCards(cfg)

	// Verify the result
	if cardTemplate == nil {
		t.Fatalf("Expected non-nil cardTemplate")
	}

	if len(cardTemplate.Cards) != len(events) {
		t.Fatalf("Expected %d cards, got %d", len(events), len(cardTemplate.Cards))
	}

	for i, card := range cardTemplate.Cards {
		if len(card.Areas) != 2 {
			t.Errorf("Expected 2 areas in card %d, got %d", i, len(card.Areas))
		}

		if len(card.Areas[1].Texts) != 2 {
			t.Errorf("Expected 2 texts in card %d, got %d", i, len(card.Areas[1].Texts))
		}

		if card.Areas[1].Texts[0].String != events[i].Desc {
			t.Errorf("Expected description '%s' in card %d, got '%s'", events[i].Desc, i, card.Areas[1].Texts[0].String)
		}

		if card.Areas[1].Texts[1].String != events[i].Title {
			t.Errorf("Expected title '%s' in card %d, got '%s'", events[i].Title, i, card.Areas[1].Texts[1].String)
		}
	}
}
