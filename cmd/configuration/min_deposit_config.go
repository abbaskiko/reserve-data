package configuration

import (
	"github.com/KyberNetwork/reserve-data/common"
)

//ExchangesMinDepositConfig store preconfig min exchange deposit for each exchange
var ExchangesMinDepositConfig = map[string]common.ExchangesMinDeposit{
	"binance": map[string]float64{
		"ETH": 0,
	},
	"bittrex": map[string]float64{
		"ETH": 0.05,
	},
	"huobi": map[string]float64{
		"ETH": 0.01,
	},
	"stable_exchange": map[string]float64{
		"ETH": 0,
	},
}
