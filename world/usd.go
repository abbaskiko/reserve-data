package world

import (
	"github.com/KyberNetwork/reserve-data/common"
)

func (tw *TheWorld) GetUSDInfo() (common.USDData, error) {
	return common.USDData{
		Timestamp:    0,
		CoinbaseUSD:  tw.getCoinbaseInfo(tw.endpoint.CoinbaseUSDEndpoint()),
		GeminiUSD:    tw.getGeminiInfo(tw.endpoint.GeminiUSDEndpoint()),
		CoinbaseUSDC: tw.getCoinbaseInfo(tw.endpoint.CoinbaseUSDCEndpoint()),
		BinanceUSDC:  tw.getBinanceInfo(tw.endpoint.BinanceUSDCEndpoint()),
		CoinbaseDAI:  tw.getCoinbaseInfo(tw.endpoint.CoinbaseDAIEndpoint()),
		HitDAI:       tw.getHitInfo(tw.endpoint.HitDaiEndpoint()),
	}, nil
}
