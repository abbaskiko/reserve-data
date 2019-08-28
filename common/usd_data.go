package common

// USDData data return by /usd-feed
type USDData struct {
	Timestamp uint64
	Coinbase  CoinbaseData `json:"CoinbaseUSD"`
	Binance   BinanceData  `json:"BinanceUSD"`
}

// BinanceData data return by https://api.binance.com/api/v3/ticker/bookTicker?symbol=ETHUSDC
type BinanceData struct {
	Valid    bool
	Error    string
	Symbol   string `json:"symbol"`
	BidPrice string `json:"bidPrice"`
	BidQty   string `json:"bidQty"`
	AskPrice string `json:"askPrice"`
	AskQty   string `json:"askQty"`
}
