package world

import (
	"fmt"

	"github.com/KyberNetwork/reserve-data/common"
)

func (tw *TheWorld) getHitInfo(url string) common.HitData {
	var result common.HitData
	err := tw.getPublic(url, &result)
	if err != nil {
		result.Error = err.Error()
		result.Valid = false
	} else {
		result.Valid = true
	}
	return result
}

func (tw *TheWorld) getBitFinexInfo(url string) common.BitFinexData {
	var result common.BitFinexData
	var bitFinexResp []float64
	err := tw.getPublic(url, &bitFinexResp)
	if err != nil {
		result.Error = err.Error()
		result.Valid = false
		return result
	}
	var bitFinexSampleResponse = []float64{201.01, 777.63092538, 201.02, 1648.5772469599997,
		-8.21238498, -0.0393, 201.01098627, 75575.33225073, 211.44, 199}
	if len(bitFinexResp) != len(bitFinexSampleResponse) {
		result.Error = fmt.Sprintf("bitfinex return unexpected number of fields, %v", bitFinexResp)
		result.Valid = false
		return result
	}

	result.Bid = bitFinexResp[0]
	result.BidSize = bitFinexResp[1]
	result.Ask = bitFinexResp[2]
	result.AskSize = bitFinexResp[3]
	result.DailyChange = bitFinexResp[4]
	result.DailyChangePerc = bitFinexResp[5]
	result.LastPrice = bitFinexResp[6]
	result.Volume = bitFinexResp[7]
	result.High = bitFinexResp[8]
	result.Low = bitFinexResp[9]

	result.Valid = true
	return result
}

// GetUSDInfo return usd info
func (tw *TheWorld) GetUSDInfo() (common.USDData, error) {
	return common.USDData{
		Timestamp:        0,
		CoinbaseUSD:      tw.getCoinbaseInfo(tw.endpoint.CoinbaseETHUSD()),
		GeminiUSD:        tw.getGeminiGoldInfo(),
		CoinbaseUSDC:     tw.getCoinbaseInfo(tw.endpoint.CoinbaseETHUSDC()),
		BinanceUSDC:      tw.getBinanceInfo(tw.endpoint.BinanceETHUSDC()),
		CoinbaseDAI:      tw.getCoinbaseInfo(tw.endpoint.CoinbaseETHDAI()),
		CoinbaseDAI10000: tw.getCoinbaseInfo(tw.endpoint.CoinbaseETHDAI10000()),
		HitDAI:           tw.getHitInfo(tw.endpoint.HitBTCETHDAI()),
		BitFinex:         tw.getBitFinexInfo(tw.endpoint.BitFinexETHUSDT()),
		BinancePAX:       tw.getBinanceInfo(tw.endpoint.BinanceETHPAX()),
		BinanceTUSD:      tw.getBinanceInfo(tw.endpoint.BinanceETHTUSD()),
		BinanceUSDT:      tw.getBinanceInfo(tw.endpoint.BinanceETHUSDT()),
	}, nil
}
