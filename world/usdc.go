package world

import (
	"github.com/KyberNetwork/reserve-data/common"
)

func (tw *TheWorld) GetUSDCInfo() (common.USDCData, error) {
	return common.USDCData{
		Timestamp: 0,
		Coinbase:  tw.getCoinbaseInfo(tw.endpoint.CoinbaseUSDCEndpoint()),
		Binance:   tw.getBinanceInfo(tw.endpoint.BinanceUSDCEndpoint()),
	}, nil
}
