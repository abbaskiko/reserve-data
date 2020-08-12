package world

import (
	"github.com/KyberNetwork/reserve-data/common"
)

// GetUSDInfo return usd info
func (tw *TheWorld) GetUSDInfo() (common.USDData, error) {
	return common.USDData{
		Timestamp:             0,
		CoinbaseETHUSDDAI5000: tw.getFeedProviderInfo(tw.endpoint.CoinbaseETHUSDDAI5000.URL),
	}, nil
}
