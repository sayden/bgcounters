package counters

type Quotes []Quote

type Quote struct {
	Origin string `json:"origin"`
	Quote  string `json:"quote"`
}
