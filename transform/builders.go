package transform

import (
	"github.com/sayden/counters"
)

// DecorateBuilder is used to compose setups of counters.Counterbuilder to execute sequentially
// 'first' will be executed first in the pipelines, then the output of 'first' will be passed
// to 'second'
func DecorateBuilder(first, second counters.CounterBuilder) counters.CounterBuilder {
	return &counterBuilderDecorator{
		first:  first,
		second: second,
	}
}

type counterBuilderDecorator struct {
	first  counters.CounterBuilder
	second counters.CounterBuilder
}

func (c *counterBuilderDecorator) ToNewCounter(counter *counters.Counter) (*counters.Counter, error) {
	newCounter, err := c.first.ToNewCounter(counter)
	if err != nil {
		return nil, err
	}

	return c.second.ToNewCounter(newCounter)
}
