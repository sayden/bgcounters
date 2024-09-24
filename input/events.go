package input

import (
	"github.com/pkg/errors"
	"github.com/sayden/counters"
	"github.com/sayden/counters/fsops"
)

func JsonEventsToEvents(eventsPoolFilepath string) ([]counters.Event, error) {
	var events []counters.Event
	if err := fsops.ReadMarkupFile(eventsPoolFilepath, &events); err != nil {
		return nil, errors.Wrap(err, "could not read events pool file")
	}

	return events, nil
}
