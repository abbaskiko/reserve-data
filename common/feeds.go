package common

type FeedConfiguration struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

type FeedSetting struct {
	BaseVolatilitySpread float64 `json:"base_volatility_spread"`
}

type MapFeedSetting map[string]FeedSetting
