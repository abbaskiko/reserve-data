package world

import (
	"encoding/json"
	"io/ioutil"
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
	OneForgeKey string `json:"oneforge"`
}

// SimulatedEndpoint for test
type SimulatedEndpoint struct {
}

// NewRealEndpointFromFile real endpoint from file
func NewRealEndpointFromFile(path string) (*RealEndpoint, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	result := RealEndpoint{}
	err = json.Unmarshal(data, &result)
	return &result, err
}

// GoldDataEndpoint real endpoint for gold
func (re RealEndpoint) GoldDataEndpoint() string {
	return "https://datafeed.digix.global/tick/"
}

// OneForgeGoldETHDataEndpoint real OneForge endpoint for gold-eth
func (re RealEndpoint) OneForgeGoldETHDataEndpoint() string {
	return "https://api.1forge.com/convert?from=XAU&to=ETH&quantity=1&api_key=" + re.OneForgeKey
}

// OneForgeGoldUSDDataEndpoint real OneForege endpoint for gold-usd
func (re RealEndpoint) OneForgeGoldUSDDataEndpoint() string {
	return "https://api.1forge.com/convert?from=XAU&to=USD&quantity=1&api_key=" + re.OneForgeKey
}

// GDAXDataEndpoint real endpoint for gdax for eht-usd
func (re RealEndpoint) GDAXDataEndpoint() string {
	return "https://api.pro.coinbase.com/products/eth-usd/ticker"
}

// KrakenDataEndpoint real kraken endpoint for eth-usd
func (re RealEndpoint) KrakenDataEndpoint() string {
	return "https://api.kraken.com/0/public/Ticker?pair=ETHUSD"
}

// GeminiDataEndpoint real gemini endpoint for eth-usd
func (re RealEndpoint) GeminiDataEndpoint() string {
	return "https://api.gemini.com/v1/pubticker/ethusd"
}

// CoinbaseBTCEndpoint real coinbase endpoint for eth-btc
func (re RealEndpoint) CoinbaseBTCEndpoint() string {
	return "https://api.pro.coinbase.com/products/eth-btc/ticker"
}

// GeminiBTCEndpoint real gemini endpoint for eth-btc
func (re RealEndpoint) GeminiBTCEndpoint() string {
	return "https://api.gemini.com/v1/pubticker/ethbtc"
}

// CoinbaseDAIEndpoint real endpoint fo Coinbase Dai
func (re RealEndpoint) CoinbaseDAIEndpoint() string {
	return "https://api.pro.coinbase.com/products/eth-dai/ticker"
}

// HitDaiEndpoint real endpoint for Hit DAI
func (re RealEndpoint) HitDaiEndpoint() string {
	return "https://api.hitbtc.com/api/2/public/ticker/ETHDAI"
}

// CoinbaseUSDEndpoint real endpoint for Coinbase USD
func (re RealEndpoint) CoinbaseUSDEndpoint() string {
	return "https://api.pro.coinbase.com/products/eth-usd/ticker"
}

// CoinbaseUSDCEndpoint real endpoint Coinbase USDC
func (re RealEndpoint) CoinbaseUSDCEndpoint() string {
	return "https://api.pro.coinbase.com/products/eth-usdc/ticker"
}

// BinanceUSDCEndpoint real endpoint
func (re RealEndpoint) BinanceUSDCEndpoint() string {
	return "https://api.binance.com/api/v3/ticker/bookTicker?symbol=ETHUSDC"
}

// BinanceUSDTEndpoint real endpoint for Binance USDT endpoint
func (re RealEndpoint) BinanceUSDTEndpoint() string {
	return "https://api.binance.com/api/v3/ticker/bookTicker?symbol=ETHUSDT"
}

// BinancePAXEndpoint real endpoint for Binance PAX
func (re RealEndpoint) BinancePAXEndpoint() string {
	return "https://api.binance.com/api/v3/ticker/bookTicker?symbol=ETHPAX"
}

// BinanceTUSDEndpoint real endpoint for binance TUSD
func (re RealEndpoint) BinanceTUSDEndpoint() string {
	return "https://api.binance.com/api/v3/ticker/bookTicker?symbol=ETHTUSD"
}

// BitFinexUSDTEndpoint real endpoint for Bitfinex USDT
func (re RealEndpoint) BitFinexUSDTEndpoint() string {
	return "https://api-pub.bitfinex.com/v2/ticker/tETHUSD"
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
