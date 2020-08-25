package configuration

import (
	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/data/fetcher"
	"github.com/KyberNetwork/reserve-data/exchange"
)

var (
	// just make sure we implement the interface
	_ fetcher.Exchange    = &exchange.Binance{}
	_ fetcher.Exchange    = &exchange.Huobi{}
	_ common.Exchange     = &exchange.Binance{}
	_ common.Exchange     = &exchange.Huobi{}
	_ common.LiveExchange = &exchange.Binance{}
	_ common.LiveExchange = &exchange.Huobi{}
)
