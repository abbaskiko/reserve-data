package http

import (
	"gopkg.in/go-playground/validator.v9"

	"github.com/KyberNetwork/reserve-data/v3/common"
)

var (
	validate = validator.New()
)

func init() {
	validate.RegisterStructValidation(createAssetInputValidation, common.CreatePendingAssetEntry{})
	validate.RegisterStructValidation(tradingPairValidation, common.TradingPair{})
}

func tradingPairValidation(sl validator.StructLevel) {
	r := sl.Current().Interface().(common.TradingPair)
	if r.Base == 0 && r.Quote == 0 {
		sl.ReportError(r.Base, "base+quote", "base and quote", "base and quote", "only one of base and quote should set")
	}
	if r.Base != 0 && r.Quote != 0 {
		sl.ReportError(r.Base, "base+quote", "base and quote", "base and quote", "only one of base and quote should set")
	}
}

func createAssetInputValidation(sl validator.StructLevel) {
	r := sl.Current().Interface().(common.CreatePendingAssetEntry)
	if !r.IsQuote && common.IsZeroAddress(r.Address) {
		sl.ReportError(r.Address, "address", "address", "address is required for non-quote", "")
	}
	if r.SetRate != common.SetRateNotSet && r.PWI == nil {
		sl.ReportError(r.PWI, "PWI", "pwi", "pwi is required for setRate set", "")
	}
	if r.Rebalance && r.RebalanceQuadratic == nil {
		sl.ReportError(r.RebalanceQuadratic, "RebalanceQuadratic", "rebalance_quadratic", "rebalance_quadratic required with rebalance", "")
	}
	if r.Rebalance && r.Target == nil {
		sl.ReportError(r.Target, "Target", "target", "target required for rebalance", "")
	}
}
