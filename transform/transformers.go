package transform

import (
	"github.com/sayden/counters"
)

// DecorateTransformer is used to compose setups of counters.Counterbuilder to execute sequentially
// 'first' will be executed first in the pipelines, then the output of 'first' will be passed
// to 'second'
func DecorateTransformer(first, second counters.CounterTransfomer) counters.CounterTransfomer {
	return &counterTransformerDecorator{
		first:  first,
		second: second,
	}
}

type counterTransformerDecorator struct {
	first  counters.CounterTransfomer
	second counters.CounterTransfomer
}

func (c *counterTransformerDecorator) ToNewCounter(counter *counters.Counter) (*counters.Counter, error) {
	newCounter, err := c.first.ToNewCounter(counter)
	if err != nil {
		return nil, err
	}

	return c.second.ToNewCounter(newCounter)
}
