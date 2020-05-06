package common

import (
	"github.com/KyberNetwork/reserve-data/common/feed"
)

// BTCData is the data returned by /btc-feed API.
type BTCData struct {
	Timestamp uint64
	Coinbase  FeedProviderResponse `json:"CoinbaseBTC"`
	Binance   FeedProviderResponse `json:"BinanceBTC"`
}

// ToMap convert to map result.
func (b BTCData) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"Timestamp":                   b.Timestamp,
		feed.CoinbaseETHBTC3.String(): b.Coinbase,
		feed.BinanceETHBTC3.String():  b.Binance,
	}
}

// CoinbaseData is the response of Coinbase ETH/BTC ticker request.
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

// GeminiDataVolume gemini data volume
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
