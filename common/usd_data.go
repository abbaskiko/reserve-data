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
	Timestamp        uint64
	CoinbaseUSD      CoinbaseData   `json:"CoinbaseUSD"`
	GeminiUSD        GeminiGoldData `json:"GeminiUSD"` // gold and usd use the same url
	CoinbaseUSDC     CoinbaseData   `json:"CoinbaseUSDC"`
	BinanceUSDC      BinanceData    `json:"BinanceUSDC"`
	CoinbaseDAI      CoinbaseData   `json:"CoinbaseDAI"`
	CoinbaseDAI10000 CoinbaseData   `json:"CoinbaseDAI10000"`
	HitDAI           HitData        `json:"HitDAI"`
	BitFinex         BitFinexData   `json:"BitFinexUSD"`
	BinanceUSDT      BinanceData    `json:"BinanceUSDT"`
	BinancePAX       BinanceData    `json:"BinancePAX"`
	BinanceTUSD      BinanceData    `json:"BinanceTUSD"`
}
