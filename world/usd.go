package world

import (
	"github.com/KyberNetwork/reserve-data/common"
)

// GetUSDInfo return usd info
func (tw *TheWorld) GetUSDInfo() (common.USDData, error) {
	return common.USDData{
		Timestamp:           0,
		CoinbaseETHDAI10000: tw.getFeedProviderInfo(tw.endpoint.CoinbaseETHDAI10000.URL),
		KrakenETHDAI10000:   tw.getFeedProviderInfo(tw.endpoint.KrakenETHDAI10000.URL),
	}, nil
}
