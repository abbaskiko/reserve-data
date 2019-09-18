package world

import (
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

func (tw *TheWorld) getCoinbaseInfo(ep string) common.CoinbaseData {
	var (
		url    = ep
		result = common.CoinbaseData{}
	)

	err := tw.getPublic(url, &result)
	if err != nil {
		result.Error = err.Error()
		result.Valid = false
	} else {
		result.Valid = true
	}
	return result
}

func (tw *TheWorld) getGeminiInfo(url string) common.GeminiData {
	var (
		result = common.GeminiData{}
	)

	err := tw.getPublic(url, &result)
	if err != nil {
		result.Error = err.Error()
		result.Valid = false
	} else {
		result.Valid = true
	}
	return result
}

func (tw *TheWorld) GetBTCInfo() (common.BTCData, error) {
	return common.BTCData{
		Coinbase: tw.getCoinbaseInfo(tw.endpoint.CoinbaseBTCEndpoint()),
		Gemini:   tw.getGeminiInfo(tw.endpoint.GeminiBTCEndpoint()),
	}, nil
}
