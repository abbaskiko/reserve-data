package common

// FeedConfigurationRequest request
type FeedConfigurationRequest struct {
	Data struct {
		Name                 string  `json:"name"`
		Enabled              bool    `json:"enabled"`
		BaseVolatilitySpread float64 `json:"base_volatility_spread"`
	} `json:"data"`
}

// FeedConfiguration is type for feed
type FeedConfiguration struct {
	Name                 string  `json:"name" db:"name"`
	Enabled              bool    `json:"enabled" db:"enabled"`
	BaseVolatilitySpread float64 `json:"base_volatility_spread" db:"base_volatility_spread"`
}
