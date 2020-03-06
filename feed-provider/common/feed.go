package common

type Feed struct {
	Ask   float64 `json:"ask,string,omitempty"`
	Bid   float64 `json:"bid,string,omitempty"`
	Error string  `json:"error,omitempty"`
	Valid bool    `json:"valid"`
}
