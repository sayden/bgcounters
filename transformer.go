package counters

type CardTransformer interface {
	ToNewCard(c *Card, sourceTemplate *CardsTemplate) (*Card, error)
}

type CounterTransfomer interface {
	ToNewCounter(c *Counter) (*Counter, error)
}

type CounterToCardTransformer interface {
	ToNewCard(c *Counter, sourceTemplate *CardsTemplate) (*Card, error)
}
