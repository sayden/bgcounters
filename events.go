package counters

type Event struct {
	Desc        string `json:"desc"`
	Title       string `json:"title"`
	InsertQuote bool   `json:"insert_quote"`
	Image       string `json:"image"`
}
