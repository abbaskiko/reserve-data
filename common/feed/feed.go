package feed

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

	CoinbaseETHBTC3 // CoinbaseETHBTCC3
	BinanceETHBTC3  // BinanceETHBTC3

	CoinbaseETHDAI10000 // CoinbaseETHDAI10000
	KrakenETHDAI10000   // KarakenETHDAI10000
)

var (
	dummyStruct = struct{}{}
	// usdFeeds list of supported usd feeds
	usdFeeds = map[string]struct{}{
		CoinbaseETHDAI10000.String(): dummyStruct,
		KrakenETHDAI10000.String():   dummyStruct,
	}

	// btcFeeds list of supported btc feeds
	btcFeeds = map[string]struct{}{
		BinanceETHBTC3.String():  dummyStruct,
		CoinbaseETHBTC3.String(): dummyStruct,
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
