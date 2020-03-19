package common

import (
	"github.com/KyberNetwork/reserve-data/world"
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

// CoinbaseData1000 ...
type CoinbaseData1000 struct {
	Valid bool    `json:"valid"`
	Error string  `json:"error"`
	Bid   float64 `json:"bid"`
	Ask   float64 `json:"ask"`
}

// USDData ...
type USDData struct {
	Timestamp        uint64
	CoinbaseUSD      CoinbaseData     `json:"CoinbaseUSD"`
	GeminiUSD        GeminiGoldData   `json:"GeminiUSD"` // gold and usd use the same url
	CoinbaseUSDC     CoinbaseData     `json:"CoinbaseUSDC"`
	BinanceUSDC      BinanceData      `json:"BinanceUSDC"`
	CoinbaseDAI      CoinbaseData     `json:"CoinbaseDAI"`
	CoinbaseDAI10000 CoinbaseData1000 `json:"CoinbaseDAI10000"`
	HitDAI           HitData          `json:"HitDAI"`
	BitFinex         BitFinexData     `json:"BitFinexUSD"`
	BinanceUSDT      BinanceData      `json:"BinanceUSDT"`
	BinancePAX       BinanceData      `json:"BinancePAX"`
	BinanceTUSD      BinanceData      `json:"BinanceTUSD"`
}

// ToMap convert to map result.
func (u USDData) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"Timestamp":                    u.Timestamp,
		world.CoinbaseETHUSD.String():  u.CoinbaseUSD,
		world.GeminiETHUSD.String():    u.GeminiUSD,
		world.CoinbaseETHUSDC.String(): u.CoinbaseUSDC,
		world.BinanceETHUSDC.String():  u.BinanceUSDC,
		world.CoinbaseETHDAI.String():  u.CoinbaseDAI,
		world.HitBTCETHDAI.String():    u.HitDAI,
		world.BitFinexETHUSDT.String(): u.BitFinex,
		world.BinanceETHUSDT.String():  u.BinanceUSDT,
		world.BinanceETHPAX.String():   u.BinancePAX,
		world.BinanceETHTUSD.String():  u.BinanceTUSD,
	}
}
