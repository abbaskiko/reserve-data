package common

// USData data return by /usd-feed
type USDData struct {
	Timestamp uint64
	Coinbase  CoinbaseData `json:"CoinbaseUSD"`
	Gemini    GeminiData   `json:"GeminiUSD"`
}
