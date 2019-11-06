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
	}

	// USDFeeds list of supported usd feeds
	USDFeeds = []string{
		"DGX",
		"OneForgeETH",
		"OneForgeUSD",
		"GDAX",
		"Kraken",
		"Gemini",
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
}

// RealEndpoint return real endpoint
type RealEndpoint struct {
	OneForgeKey string `json:"oneforge"`
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

// SimulatedEndpoint for test
type SimulatedEndpoint struct {
}

// GoldDataEndpoint simulated endpoint
func (se SimulatedEndpoint) GoldDataEndpoint() string {
	return "http://simulator:5400/tick"
}

// OneForgeGoldETHDataEndpoint simulated endpoint
func (se SimulatedEndpoint) OneForgeGoldETHDataEndpoint() string {
	return "http://simulator:5500/1.0.3/convert?from=XAU&to=ETH&quantity=1&api_key="
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

// GeminiDataEndpoint simulated endpoint
func (se SimulatedEndpoint) GeminiDataEndpoint() string {
	return "http://simulator:5800/v1/pubticker/ethusd"
}

// CoinbaseBTCEndpoint simulated endpoint fo coinbase btc feed
func (se SimulatedEndpoint) CoinbaseBTCEndpoint() string {
	return "http://simulator:5600/products/eth-btc/ticker"
}

// GeminiBTCEndpoint gemini simulated endpoint
func (se SimulatedEndpoint) GeminiBTCEndpoint() string {
	return "http://simulator:5800/v1/pubticker/ethbtc"
}
