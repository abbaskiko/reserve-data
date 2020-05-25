package common

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

// HitData ...
type HitData struct {
	Valid       bool
	Error       string
	Ask         string `json:"ask"`
	Bid         string `json:"bid"`
	Last        string `json:"last"`
	Open        string `json:"open"`
	Low         string `json:"low"`
	High        string `json:"high"`
	Volume      string `json:"volume"`
	VolumeQuote string `json:"volumeQuote"`
	Timestamp   string `json:"timestamp"`
	Symbol      string `json:"symbol"`
}

// USDData ...
type USDData struct {
	Timestamp            uint64
	CoinbaseETHUSD10000  FeedProviderResponse `json:"CoinbaseETHUSD10000"`
	GeminiETHUSD10000    FeedProviderResponse `json:"GeminiETHUSD10000"`
	CoinbaseETHUSDC10000 FeedProviderResponse `json:"CoinbaseETHUSDC10000"`
	BinanceETHUSDC10000  FeedProviderResponse `json:"BinanceETHUSDC10000"`
	CoinbaseETHDAI5000   FeedProviderResponse `json:"CoinbaseETHDAI5000"`
	BitfinexETHUSDT10000 FeedProviderResponse `json:"BitfinexETHUSDT10000"`
	BinanceETHUSDT10000  FeedProviderResponse `json:"BinanceETHUSDT10000"`
	BinanceETHPAX5000    FeedProviderResponse `json:"BinanceETHPAX5000"`
	BinanceETHBUSD100000 FeedProviderResponse `json:"BinanceETHBUSD100000"`
}
