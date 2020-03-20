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

	CoinbaseETHBTC // CoinbaseETHBTC
	GeminiETHBTC   // GeminiETHBTC

	CoinbaseETHUSDC     // CoinbaseETHUSDC
	BinanceETHUSDC      // BinanceETHUSDC
	CoinbaseETHUSD      // CoinbaseETHUSD
	CoinbaseETHDAI      // CoinbaseETHDAI
	CoinbaseETHDAI10000 // CoinbaseETHDAI10000
	HitBTCETHDAI        // HitBTCETHDAI
	BitFinexETHUSDT     // BitFinexETHUSDT
	BinanceETHUSDT      // BinanceETHUSDT
	BinanceETHPAX       // BinanceETHPAX
	BinanceETHTUSD      // BinanceETHTUSD
)

var (
	dummyStruct = struct{}{}
	// usdFeeds list of supported usd feeds
	usdFeeds = map[string]struct{}{
		CoinbaseETHUSD.String():      dummyStruct,
		GeminiETHUSD.String():        dummyStruct,
		CoinbaseETHUSDC.String():     dummyStruct,
		BinanceETHUSDC.String():      dummyStruct,
		CoinbaseETHDAI.String():      dummyStruct,
		CoinbaseETHDAI10000.String(): dummyStruct,
		HitBTCETHDAI.String():        dummyStruct,
		BitFinexETHUSDT.String():     dummyStruct,
		BinanceETHPAX.String():       dummyStruct,
		BinanceETHTUSD.String():      dummyStruct,
		BinanceETHUSDT.String():      dummyStruct,
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
