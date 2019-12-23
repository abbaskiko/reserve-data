package common

type FeedConfiguration struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

// FeedSetting setting for feed configuration
type FeedSetting struct {
	BaseVolatilitySpread float64 `json:"base_volatility_spread"`
	NormalSpread         float64 `json:"normal_spread"`
}

// MapFeedSetting map feed name with feed setting
type MapFeedSetting map[string]FeedSetting
