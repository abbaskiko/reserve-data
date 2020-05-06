package common

import (
	"github.com/KyberNetwork/reserve-data/common/feed"
)

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

// USDData ...
type USDData struct {
	Timestamp           uint64
	CoinbaseETHDAI10000 FeedProviderResponse `json:"CoinbaseDAI10000"`
	KrakenETHDAI10000   FeedProviderResponse `json:"KrakenETHDAI10000"`
}

// ToMap convert to map result.
func (u USDData) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"Timestamp":                       u.Timestamp,
		feed.CoinbaseETHDAI10000.String(): u.CoinbaseETHDAI10000,
		feed.KrakenETHDAI10000.String():   u.KrakenETHDAI10000,
	}
}
