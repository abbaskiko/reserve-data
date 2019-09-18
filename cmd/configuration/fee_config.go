package configuration

import "github.com/KyberNetwork/reserve-data/common"

//FeeConfigs store predefined fee configs of exchanges
var FeeConfigs = map[string]common.ExchangeFees{
	"binance": {
		Trading: common.TradingFee{
			"taker": 0.001,
			"maker": 0.001,
		},
		Funding: common.FundingFee{
			Withdraw: map[string]float64{
				"ETH": 0.01,
			},
			Deposit: map[string]float64{
				"ETH": 0,
			},
		},
	},
	"bittrex": {
		Trading: common.TradingFee{
			"taker": 0.0025,
			"maker": 0.0025,
		},
		Funding: common.FundingFee{
			Withdraw: map[string]float64{
				"ETH": 0.006,
			},
			Deposit: map[string]float64{
				"ETH": 0,
			},
		},
	},
	"huobi": {
		Trading: common.TradingFee{
			"taker": 0.002,
			"maker": 0.002,
		},
		Funding: common.FundingFee{
			Withdraw: map[string]float64{
				"ETH": 0.01,
			},
			Deposit: map[string]float64{
				"ETH": 0.01,
			},
		},
	},
	"stable_exchange": {
		Trading: common.TradingFee{
			"taker": 0,
			"maker": 0,
		},
		Funding: common.FundingFee{
			Withdraw: map[string]float64{
				"ETH": 0,
			},
			Deposit: map[string]float64{
				"ETH": 0,
			},
		},
	},
}
