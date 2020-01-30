package world

import (
	"encoding/json"
	"io/ioutil"

	"github.com/KyberNetwork/reserve-data/common"
)

var (
	allFeeds = []string{
		"DGX",
		"OneForgeETH",
		"OneForgeUSD",
		"GDAX",
		"Kraken",
		"Gemini",

		"CoinbaseBTC",
		"GeminiBTC",

		"CoinbaseUSD",

		"CoinbaseUSDC",
		"BinanceUSDC",

		"CoinbaseDAI",
		"HitDAI",

		"BitFinexUSD",
		"BinanceUSDT",
		"BinancePAX",
		"BinanceTUSD",
	}

	// USDFeeds list of supported usd feeds
	USDFeeds = []string{
		"OneForgeUSD",
		"CoinbaseUSD",
		"BinanceUSD",

		"CoinbaseUSDC",
		"BinanceUSDC",

		"CoinbaseDAI",
		"HitDAI",

		"BitFinexUSD",
		"BinanceUSDT",
		"BinancePAX",
		"BinanceTUSD",
	}

	// BTCFeeds list of supported btc feeds
	BTCFeeds = []string{
		"CoinbaseBTC",
		"GeminiBTC",
	}
)

// AllFeeds returns all configured feed sources.
func AllFeeds() []string {
	return allFeeds
}

// Endpoint returns all API endpoints to use in TheWorld struct.
type Endpoint interface {
	GoldDataEndpoint() string
	OneForgeGoldETHDataEndpoint() string
	OneForgeGoldUSDDataEndpoint() string
	GDAXDataEndpoint() string
	KrakenDataEndpoint() string
	GeminiDataEndpoint() string

	CoinbaseBTCEndpoint() string
	GeminiBTCEndpoint() string

	CoinbaseUSDCEndpoint() string
	BinanceUSDCEndpoint() string
	CoinbaseUSDEndpoint() string
	CoinbaseDAIEndpoint() string
	HitDaiEndpoint() string

	BitFinexUSDTEndpoint() string
	BinanceUSDTEndpoint() string
	BinancePAXEndpoint() string
	BinanceTUSDEndpoint() string
}

// RealEndpoint return real endpoint
type RealEndpoint struct {
	Endpoints common.WorldEndpoints `json:"endpoints"`
}

// SimulatedEndpoint for test
type SimulatedEndpoint struct {
}

// NewRealEndpointFromFile real endpoint from file
func NewRealEndpointFromFile(path string) (*RealEndpoint, error) {
	result := &RealEndpoint{}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(data, result); err != nil {
		return nil, err
	}
	return result, nil
}

// GoldDataEndpoint real endpoint for gold
func (re RealEndpoint) GoldDataEndpoint() string {
	return re.Endpoints.GoldData.URL
}

// OneForgeGoldETHDataEndpoint real OneForge endpoint for gold-eth
func (re RealEndpoint) OneForgeGoldETHDataEndpoint() string {
	return re.Endpoints.OneForgeGoldETH.URL
}

// OneForgeGoldUSDDataEndpoint real OneForege endpoint for gold-usd
func (re RealEndpoint) OneForgeGoldUSDDataEndpoint() string {
	return re.Endpoints.OneForgeGoldUSD.URL
}

// GDAXDataEndpoint real endpoint for gdax for eht-usd
func (re RealEndpoint) GDAXDataEndpoint() string {
	return re.Endpoints.GDAXData.URL
}

// KrakenDataEndpoint real kraken endpoint for eth-usd
func (re RealEndpoint) KrakenDataEndpoint() string {
	return re.Endpoints.KrakenData.URL
}

// GeminiDataEndpoint real gemini endpoint for eth-usd
func (re RealEndpoint) GeminiDataEndpoint() string {
	return re.Endpoints.GeminiData.URL
}

// CoinbaseBTCEndpoint real coinbase endpoint for eth-btc
func (re RealEndpoint) CoinbaseBTCEndpoint() string {
	return re.Endpoints.CoinbaseBTC.URL
}

// GeminiBTCEndpoint real gemini endpoint for eth-btc
func (re RealEndpoint) GeminiBTCEndpoint() string {
	return re.Endpoints.GeminiBTC.URL
}

// CoinbaseDAIEndpoint real endpoint fo Coinbase Dai
func (re RealEndpoint) CoinbaseDAIEndpoint() string {
	return re.Endpoints.CoinbaseDAI.URL
}

// HitDaiEndpoint real endpoint for Hit DAI
func (re RealEndpoint) HitDaiEndpoint() string {
	return re.Endpoints.HitDai.URL
}

// CoinbaseUSDEndpoint real endpoint for Coinbase USD
func (re RealEndpoint) CoinbaseUSDEndpoint() string {
	return re.Endpoints.CoinbaseUSD.URL
}

// CoinbaseUSDCEndpoint real endpoint Coinbase USDC
func (re RealEndpoint) CoinbaseUSDCEndpoint() string {
	return re.Endpoints.CoinbaseUSDC.URL
}

// BinanceUSDCEndpoint real endpoint
func (re RealEndpoint) BinanceUSDCEndpoint() string {
	return re.Endpoints.BinanceUSDC.URL
}

// BinanceUSDTEndpoint real endpoint for Binance USDT endpoint
func (re RealEndpoint) BinanceUSDTEndpoint() string {
	return re.Endpoints.BinanceUSDT.URL
}

// BinancePAXEndpoint real endpoint for Binance PAX
func (re RealEndpoint) BinancePAXEndpoint() string {
	return re.Endpoints.BinancePAX.URL
}

// BinanceTUSDEndpoint real endpoint for binance TUSD
func (re RealEndpoint) BinanceTUSDEndpoint() string {
	return re.Endpoints.BinanceTUSD.URL
}

// BitFinexUSDTEndpoint real endpoint for Bitfinex USDT
func (re RealEndpoint) BitFinexUSDTEndpoint() string {
	return re.Endpoints.BitFinexUSDT.URL
}
