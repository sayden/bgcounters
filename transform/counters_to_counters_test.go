package transform

import (
	"testing"

	"github.com/sayden/counters"
	"github.com/stretchr/testify/assert"
)

func TestCountersToCounters(t *testing.T) {
	originalTemplate := &counters.CounterTemplate{
		Counters: []counters.Counter{
			{Texts: counters.Texts{{String: "name_transformed"}}},
		}}

	cfg := &CountersToCountersConfig{
		OriginalCounterTemplate: originalTemplate,
		OutputPathInTemplate:    "output/path",
		CounterTransformer:      &SimpleFowCounterBuilder{}}

	result, err := cfg.CountersToCounters()
	assert.NoError(t, err)
	assert.Equal(t, "output/path", result.OutputFolder)
	assert.Equal(t, 1, len(result.Counters))
	assert.Equal(t, 0, len(result.Counters[0].Texts))
}
