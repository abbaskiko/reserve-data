package common

// USDCData data return by /usdc-feed
type USDCData struct {
	Timestamp uint64
	Coinbase  CoinbaseData `json:"CoinbaseUSDC"`
	Binance   BinanceData  `json:"BinanceUSDC"`
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
