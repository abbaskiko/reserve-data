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
	EPS common.WorldEndpoints `json:"eps"`
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
	return re.EPS.GoldData.URL
}

// OneForgeGoldETHDataEndpoint real OneForge endpoint for gold-eth
func (re RealEndpoint) OneForgeGoldETHDataEndpoint() string {
	return re.EPS.OneForgeGoldETH.URL
}

// OneForgeGoldUSDDataEndpoint real OneForege endpoint for gold-usd
func (re RealEndpoint) OneForgeGoldUSDDataEndpoint() string {
	return re.EPS.OneForgeGoldUSD.URL
}

// GDAXDataEndpoint real endpoint for gdax for eht-usd
func (re RealEndpoint) GDAXDataEndpoint() string {
	return re.EPS.GDAXData.URL
}

// KrakenDataEndpoint real kraken endpoint for eth-usd
func (re RealEndpoint) KrakenDataEndpoint() string {
	return re.EPS.KrakenData.URL
}

// GeminiDataEndpoint real gemini endpoint for eth-usd
func (re RealEndpoint) GeminiDataEndpoint() string {
	return re.EPS.GeminiData.URL
}

// CoinbaseBTCEndpoint real coinbase endpoint for eth-btc
func (re RealEndpoint) CoinbaseBTCEndpoint() string {
	return re.EPS.CoinbaseBTC.URL
}

// GeminiBTCEndpoint real gemini endpoint for eth-btc
func (re RealEndpoint) GeminiBTCEndpoint() string {
	return re.EPS.GeminiBTC.URL
}

// CoinbaseDAIEndpoint real endpoint fo Coinbase Dai
func (re RealEndpoint) CoinbaseDAIEndpoint() string {
	return re.EPS.CoinbaseDAI.URL
}

// HitDaiEndpoint real endpoint for Hit DAI
func (re RealEndpoint) HitDaiEndpoint() string {
	return re.EPS.HitDai.URL
}

// CoinbaseUSDEndpoint real endpoint for Coinbase USD
func (re RealEndpoint) CoinbaseUSDEndpoint() string {
	return re.EPS.CoinbaseUSD.URL
}

// CoinbaseUSDCEndpoint real endpoint Coinbase USDC
func (re RealEndpoint) CoinbaseUSDCEndpoint() string {
	return re.EPS.CoinbaseUSDC.URL
}

// BinanceUSDCEndpoint real endpoint
func (re RealEndpoint) BinanceUSDCEndpoint() string {
	return re.EPS.BinanceUSDC.URL
}

// BinanceUSDTEndpoint real endpoint for Binance USDT endpoint
func (re RealEndpoint) BinanceUSDTEndpoint() string {
	return re.EPS.BinanceUSDT.URL
}

// BinancePAXEndpoint real endpoint for Binance PAX
func (re RealEndpoint) BinancePAXEndpoint() string {
	return re.EPS.BinancePAX.URL
}

// BinanceTUSDEndpoint real endpoint for binance TUSD
func (re RealEndpoint) BinanceTUSDEndpoint() string {
	return re.EPS.BinanceTUSD.URL
}

// BitFinexUSDTEndpoint real endpoint for Bitfinex USDT
func (re RealEndpoint) BitFinexUSDTEndpoint() string {
	return re.EPS.BitFinexUSDT.URL
}

// GeminiDataEndpoint simulated endpoint
func (se SimulatedEndpoint) GeminiDataEndpoint() string {
	return "http://simulator:5800/v1/pubticker/ethusd"
}

// OneForgeGoldUSDDataEndpoint simulated endpoint
func (se SimulatedEndpoint) OneForgeGoldUSDDataEndpoint() string {
	return "http://simulator:5500/1.0.3/convert?from=XAU&to=USD&quantity=1&api_key="
}

// GDAXDataEndpoint simulated endpoint
func (se SimulatedEndpoint) GDAXDataEndpoint() string {
	return "http://simulator:5600/products/eth-usd/ticker"
}

// KrakenDataEndpoint simulated endpoint
func (se SimulatedEndpoint) KrakenDataEndpoint() string {
	return "http://simulator:5700/0/public/Ticker?pair=ETHUSD"
}

// GoldDataEndpoint simulated endpoint
func (se SimulatedEndpoint) GoldDataEndpoint() string {
	return "http://simulator:5400/tick"
}

// OneForgeGoldETHDataEndpoint simulated endpoint
func (se SimulatedEndpoint) OneForgeGoldETHDataEndpoint() string {
	return "http://simulator:5500/1.0.3/convert?from=XAU&to=ETH&quantity=1&api_key="
}

// GeminiBTCEndpoint gemini simulated endpoint
func (se SimulatedEndpoint) GeminiBTCEndpoint() string {
	return "http://simulator:5800/v1/pubticker/ethbtc"
}

// CoinbaseBTCEndpoint simulated endpoint fo coinbase btc feed
func (se SimulatedEndpoint) CoinbaseBTCEndpoint() string {
	return "http://simulator:5600/products/eth-btc/ticker"
}

// TODO: support simulation

// CoinbaseUSDCEndpoint simulator endpoint for Coinbase USDC
func (se SimulatedEndpoint) CoinbaseUSDCEndpoint() string {
	panic("implement me")
}

// BinanceUSDCEndpoint simulator endpoint for Binance USDC
func (se SimulatedEndpoint) BinanceUSDCEndpoint() string {
	panic("implement me")
}

// BinancePAXEndpoint simulator endpoint for Binance PAX
func (se SimulatedEndpoint) BinancePAXEndpoint() string {
	panic("implement me")
}

// BinanceUSDTEndpoint simulator endpoint for Binance USDT
func (se SimulatedEndpoint) BinanceUSDTEndpoint() string {
	panic("implement me")
}

// BitFinexUSDTEndpoint simulator endpoint for Bitfinex USDT
func (se SimulatedEndpoint) BitFinexUSDTEndpoint() string {
	panic("implement me")
}

// BinanceTUSDEndpoint simulator endpoint for Binance TUSD
func (se SimulatedEndpoint) BinanceTUSDEndpoint() string {
	panic("implement me")
}

// CoinbaseDAIEndpoint simulator endpoint for Coinbase DAI
func (se SimulatedEndpoint) CoinbaseDAIEndpoint() string {
	panic("implement me")
}

// HitDaiEndpoint simulator endpoint for Hit Dai
func (se SimulatedEndpoint) HitDaiEndpoint() string {
	panic("implement me")
}

// CoinbaseUSDEndpoint simulator endpoint for Coinbase USD
func (se SimulatedEndpoint) CoinbaseUSDEndpoint() string {
	panic("implement me")
}
