package transform

import (
	"testing"

	"github.com/sayden/counters"
	"github.com/stretchr/testify/assert"
)

type mockCounterBuilder struct {
	counter *counters.Counter
	err     error
}

func (m *mockCounterBuilder) ToNewCounter(counter *counters.Counter) (*counters.Counter, error) {
	if m.err != nil {
		return nil, m.err
	}
	counter.Texts[0].String += " world"
	return m.counter, nil
}

func TestDecorateBuilder(t *testing.T) {
	counter1 := &counters.Counter{Texts: counters.Texts{{String: "hello"}}}
	counter2 := &counters.Counter{Texts: counters.Texts{{String: "hello"}}}
	counter3 := &counters.Counter{Texts: counters.Texts{{String: "hello world"}}}

	firstBuilder := &mockCounterBuilder{counter: counter2, err: nil}
	secondBuilder := &mockCounterBuilder{counter: counter3, err: nil}

	decorator := DecorateTransformer(firstBuilder, secondBuilder)

	result, err := decorator.ToNewCounter(counter1)
	assert.NoError(t, err)
	assert.Equal(t, counter3, result)
}

func TestDecorateBuilder_FirstBuilderError(t *testing.T) {
	counter1 := &counters.Counter{Texts: counters.Texts{{String: "hello"}}}
	expectedError := assert.AnError

	firstBuilder := &mockCounterBuilder{counter: nil, err: expectedError}
	secondBuilder := &mockCounterBuilder{counter: nil, err: nil}

	decorator := DecorateTransformer(firstBuilder, secondBuilder)

	result, err := decorator.ToNewCounter(counter1)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
}

func TestDecorateBuilder_SecondBuilderError(t *testing.T) {
	counter1 := &counters.Counter{Texts: counters.Texts{{String: "hello"}}}
	counter2 := &counters.Counter{Texts: counters.Texts{{String: "hello"}}}
	expectedError := assert.AnError

	firstBuilder := &mockCounterBuilder{counter: counter2, err: nil}
	secondBuilder := &mockCounterBuilder{counter: nil, err: expectedError}

	decorator := DecorateTransformer(firstBuilder, secondBuilder)

	result, err := decorator.ToNewCounter(counter1)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
}
