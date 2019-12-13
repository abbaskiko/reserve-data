package common

// FeedConfiguration ...
type FeedConfiguration struct {
	Name                 string  `json:"name"`
	Enabled              bool    `json:"enabled"`
	BaseVolatilitySpread float64 `json:"base_volatility_spread"`
}
