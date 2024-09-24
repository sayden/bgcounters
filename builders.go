package counters

type CardBuilder interface {
	ToCard(c Counter, sourceTemplate *CardsTemplate) (*Card, error)
}

type CounterBuilder interface {
	ToCounter(c *Counter) (*Counter, error)
}
