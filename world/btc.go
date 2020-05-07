package world

import (
	"github.com/KyberNetwork/reserve-data/common"
)

// GetBTCInfo return btc info
func (tw *TheWorld) GetBTCInfo() (common.BTCData, error) {
	return common.BTCData{
		Coinbase: tw.getFeedProviderInfo(tw.endpoint.CoinbaseETHBTC3.URL),
		Binance:  tw.getFeedProviderInfo(tw.endpoint.BinanceETHBTC3.URL),
	}, nil
}
