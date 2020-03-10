package world

import (
	"github.com/KyberNetwork/reserve-data/common"
)

func (tw *TheWorld) getCoinbaseInfo(ep string) common.CoinbaseData {
	var (
		url    = ep
		result = common.CoinbaseData{
			Valid: true,
		}
	)
	if err := tw.getPublic(url, &result); err != nil {
		result.Valid = false
		result.Error = err.Error()
	}
	return result
}

func (tw *TheWorld) getGeminiInfo(endpoint string) common.GeminiData {
	var (
		url    = endpoint
		result = common.GeminiData{
			Valid: true,
		}
	)
	if err := tw.getPublic(url, &result); err != nil {
		result.Valid = false
		result.Error = err.Error()
	}
	return result
}

// GetBTCInfo return btc info
func (tw *TheWorld) GetBTCInfo() (common.BTCData, error) {
	return common.BTCData{
		Coinbase: tw.getCoinbaseInfo(tw.endpoint.CoinbaseETHBTC()),
		Gemini:   tw.getGeminiInfo(tw.endpoint.GeminiETHBTC()),
	}, nil
}
