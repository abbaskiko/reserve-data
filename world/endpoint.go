package world

import (
	"encoding/json"
	"io/ioutil"

	"github.com/KyberNetwork/reserve-data/common"
)

type Feed int

// Feed ..
//go:generate stringer -type=Feed -linecomment
const (
	GoldData       Feed = iota + 1 // GoldData
	OneForgeXAUETH                 // OneForgeXAUETH
	OneForgeXAUUSD                 // OneForgeXAUUSD
	GDAXETHUSD                     // GDAXETHUSD
	KrakenETHUSD                   // KrakenETHUSD
	GeminiETHUSD                   // GeminiETHUSD

	CoinbaseETHBTC // CoinbaseETHBTC
	GeminiETHBTC   // GeminiETHBTC

	CoinbaseETHUSDC // CoinbaseETHUSDC
	BinanceETHUSDC  // BinanceETHUSDC
	CoinbaseETHUSD  // CoinbaseETHUSD
	CoinbaseETHDAI  // CoinbaseETHDAI
	HitBTCETHDAI    // HitBTCETHDAI
	BitFinexETHUSDT // BitFinexETHUSDT
	BinanceETHUSDT  // BinanceETHUSDT
	BinanceETHPAX   // BinanceETHPAX
	BinanceETHTUSD  // BinanceETHTUSD
)

var (
	dummyStruct = struct{}{}
	// usdFeeds list of supported usd feeds
	usdFeeds = map[string]struct{}{
		CoinbaseETHUSD.String():  dummyStruct,
		GeminiETHUSD.String():    dummyStruct,
		CoinbaseETHUSDC.String(): dummyStruct,
		BinanceETHUSDC.String():  dummyStruct,
		CoinbaseETHDAI.String():  dummyStruct,
		HitBTCETHDAI.String():    dummyStruct,
		BitFinexETHUSDT.String(): dummyStruct,
		BinanceETHPAX.String():   dummyStruct,
		BinanceETHTUSD.String():  dummyStruct,
		BinanceETHUSDT.String():  dummyStruct,
	}

	// btcFeeds list of supported btc feeds
	btcFeeds = map[string]struct{}{
		CoinbaseETHBTC.String(): dummyStruct,
		GeminiETHBTC.String():   dummyStruct,
	}

	goldFeeds = map[string]struct{}{
		GoldData.String():       dummyStruct,
		OneForgeXAUETH.String(): dummyStruct,
		OneForgeXAUUSD.String(): dummyStruct,
		GDAXETHUSD.String():     dummyStruct,
		KrakenETHUSD.String():   dummyStruct,
		GeminiETHUSD.String():   dummyStruct,
	}
)

// SupportedFeed ...
type SupportedFeed struct {
	Gold map[string]struct{}
	BTC  map[string]struct{}
	USD  map[string]struct{}
}

// AllFeeds returns all configured feed sources.
func AllFeeds() SupportedFeed {
	return SupportedFeed{
		Gold: goldFeeds,
		BTC:  btcFeeds,
		USD:  usdFeeds,
	}
}

// Endpoint returns all API endpoints to use in TheWorld struct.
type Endpoint interface {
	GoldDataEndpoint() string
	OneForgeXAUETH() string
	OneForgeXAUUSD() string
	GDAXETHUSD() string
	KrakenETHUSD() string
	GeminiETHUSD() string

	CoinbaseETHBTC() string
	GeminiETHBTC() string

	CoinbaseETHUSDC() string
	BinanceETHUSDC() string
	CoinbaseETHUSD() string
	CoinbaseETHDAI() string
	HitBTCETHDAI() string

	BitFinexETHUSDT() string
	BinanceETHUSDT() string
	BinanceETHPAX() string
	BinanceETHTUSD() string
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

// OneForgeXAUETH real OneForge endpoint for gold-eth
func (re RealEndpoint) OneForgeXAUETH() string {
	return re.Endpoints.OneForgeGoldETH.URL
}

// OneForgeXAUUSD real OneForege endpoint for gold-usd
func (re RealEndpoint) OneForgeXAUUSD() string {
	return re.Endpoints.OneForgeGoldUSD.URL
}

// GDAXETHUSD real endpoint for gdax for eht-usd
func (re RealEndpoint) GDAXETHUSD() string {
	return re.Endpoints.GDAXData.URL
}

// KrakenETHUSD real kraken endpoint for eth-usd
func (re RealEndpoint) KrakenETHUSD() string {
	return re.Endpoints.KrakenData.URL
}

// GeminiETHUSD real gemini endpoint for eth-usd
func (re RealEndpoint) GeminiETHUSD() string {
	return re.Endpoints.GeminiData.URL
}

// CoinbaseETHBTC real coinbase endpoint for eth-btc
func (re RealEndpoint) CoinbaseETHBTC() string {
	return re.Endpoints.CoinbaseBTC.URL
}

// GeminiETHBTC real gemini endpoint for eth-btc
func (re RealEndpoint) GeminiETHBTC() string {
	return re.Endpoints.GeminiBTC.URL
}

// CoinbaseETHDAI real endpoint fo Coinbase Dai
func (re RealEndpoint) CoinbaseETHDAI() string {
	return re.Endpoints.CoinbaseDAI.URL
}

// HitBTCETHDAI real endpoint for Hit DAI
func (re RealEndpoint) HitBTCETHDAI() string {
	return re.Endpoints.HitDai.URL
}

// CoinbaseETHUSD real endpoint for Coinbase USD
func (re RealEndpoint) CoinbaseETHUSD() string {
	return re.Endpoints.CoinbaseUSD.URL
}

// CoinbaseETHUSDC real endpoint Coinbase USDC
func (re RealEndpoint) CoinbaseETHUSDC() string {
	return re.Endpoints.CoinbaseUSDC.URL
}

// BinanceETHUSDC real endpoint
func (re RealEndpoint) BinanceETHUSDC() string {
	return re.Endpoints.BinanceUSDC.URL
}

// BinanceETHUSDT real endpoint for Binance USDT endpoint
func (re RealEndpoint) BinanceETHUSDT() string {
	return re.Endpoints.BinanceUSDT.URL
}

// BinanceETHPAX real endpoint for Binance PAX
func (re RealEndpoint) BinanceETHPAX() string {
	return re.Endpoints.BinancePAX.URL
}

// BinanceETHTUSD real endpoint for binance TUSD
func (re RealEndpoint) BinanceETHTUSD() string {
	return re.Endpoints.BinanceTUSD.URL
}

// BitFinexETHUSDT real endpoint for Bitfinex USDT
func (re RealEndpoint) BitFinexETHUSDT() string {
	return re.Endpoints.BitFinexUSDT.URL
}
