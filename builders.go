package counters

type CardBuilder interface {
	ToNewCard(c Counter, sourceTemplate *CardsTemplate) (*Card, error)
}

type CounterBuilder interface {
	ToNewCounter(c *Counter) (*Counter, error)
}
