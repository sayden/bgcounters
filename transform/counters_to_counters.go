package transform

import (
	"github.com/pkg/errors"
	"github.com/sayden/counters"
)

type CountersToCountersConfig struct {
	OriginalCounterTemplate *counters.CounterTemplate
	OutputPathInTemplate    string
	CounterBuilder          counters.CounterBuilder
}

func CountersToCounters(cfg *CountersToCountersConfig) (*counters.CounterTemplate, error) {
	cfg.OriginalCounterTemplate.OutputFolder = cfg.OutputPathInTemplate

	finalOutputCounters := make([]counters.Counter, 0)

	for _, counter := range cfg.OriginalCounterTemplate.Counters {
		newCounter, err := cfg.CounterBuilder.ToNewCounter(&counter)
		if err != nil {
			return nil, errors.Wrap(err, "error trying to build counter")
		}

		finalOutputCounters = append(finalOutputCounters, *newCounter)
	}

	cfg.OriginalCounterTemplate.Counters = nil
	cfg.OriginalCounterTemplate.Counters = finalOutputCounters

	return cfg.OriginalCounterTemplate, nil
}
