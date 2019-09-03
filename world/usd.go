package world

import (
	"github.com/KyberNetwork/reserve-data/common"
)

func (tw *TheWorld) GetUSDInfo() (common.USDData, error) {
	return common.USDData{
		Timestamp: 0,
		Coinbase:  tw.getCoinbaseInfo(tw.endpoint.CoinbaseUSDEndpoint()),
		Binance:   tw.getBinanceInfo(tw.endpoint.BinanceUSDEndpoint()),
	}, nil
}
