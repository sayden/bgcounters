package input

import (
	"github.com/pkg/errors"
	"github.com/sayden/counters"
	"github.com/sayden/counters/fsops"
)

func ReadQuotesFromFile(quotesFilepath string) (counters.Quotes, error) {
	var qs counters.Quotes
	if err := fsops.ReadMarkupFile(quotesFilepath, &qs); err != nil {
		return nil, errors.Wrap(err, "could not read quotes file")
	}

	return qs, nil
}
