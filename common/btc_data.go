package common

// BitfinexData is the response of Bitfinix Ticker request.
// Example:
//
// Request: https://api.bitfinex.com/v1/pubticker/btcusd
// Response:
// {
// "mid":"244.755",
// "bid":"244.75",
// "ask":"244.76",
// "last_price":"244.82",
// "low":"244.2",
// "high":"248.19",
// "volume":"7842.11542563",
// "timestamp":"1444253422.348340958"
//}
type BitfinexData struct {
	Valid     bool
	Error     string
	Mid       string `json:"mid"`
	Bid       string `json:"bid"`
	Ask       string `json:"ask"`
	LastPrice string `json:"last_price"`
	Low       string `json:"low"`
	High      string `json:"high"`
	Volume    string `json:"volume"`
	Timestamp string `json:"timestamp"`
}

// BinanceData is the response of Binance Ticker request.
// Example:
// Request: https://api.binance.com/api/v3/ticker/bookTicker?symbol=ETHBTC
// Response: {
//  "symbol": "ETHBTC",
//  "bidPrice": "0.03338700",
//  "bidQty": "3.39200000",
//  "askPrice": "0.03339400",
//  "askQty": "0.08600000"
//}
type BinanceData struct {
	Valid    bool
	Error    string
	Symbol   string `json:"symbol"`
	BidPrice string `json:"bidPrice"`
	BidQty   string `json:"bidQty"`
	AskPrice string `json:"askPrice"`
	AskQty   string `json:"askQty"`
}

// BTCData is the data returned by /btc-feed API.
type BTCData struct {
	Timestamp uint64
	Bitfinex  BitfinexData `json:"bitfinex"`
	Binance   BinanceData  `json:"binance"`
	Coinbase  CoinbaseData `json:"coinbase"`
	Gemini    GeminiData   `json:"gemini"`
}

// CoinBaseData is the response of Coinbase ETH/BTC ticker request.
// Example: https://api.pro.coinbase.com/products/eth-btc/ticker
// Response
// {
//  "trade_id": 7188449,
//  "price": "0.02707000",
//  "size": "3.43340266",
//  "time": "2019-05-13T04:56:14.777Z",
//  "bid": "0.02707",
//  "ask": "0.02708",
//  "volume": "31528.93120442"
//}
type CoinbaseData struct {
	Valid   bool
	Error   string
	TradeID uint64 `json:"trade_id"`
	Price   string `json:"price"`
	Size    string `json:"size"`
	Time    string `json:"time"`
	Bid     string `json:"bid"`
	Ask     string `json:"ask"`
	Volume  string `json:"volume"`
}

type GeminiDataVolume struct {
	ETH       string `json:"ETH"`
	BTC       string `json:"BTC"`
	Timestamp uint64 `json:"timestamp"`
}

// GeminiData is the data returns by Gemini ETH/BTC ticker.
// Example: https://api.gemini.com/v1/pubticker/ethbtc
// Response
// {
//  "bid": "0.02686",
//  "ask": "0.02694",
//  "volume": {
//    "ETH": "5595.131076",
//    "BTC": "149.60020830467",
//    "timestamp": 1557728700000
//  },
//  "last": "0.02694"
// }
type GeminiData struct {
	Valid  bool
	Error  string
	Bid    string           `json:"bid"`
	Ask    string           `json:"ask"`
	Volume GeminiDataVolume `json:"volume"`
	Last   string           `json:"last"`
}
