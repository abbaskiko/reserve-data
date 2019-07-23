package http

import (
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v8"

	v3 "github.com/KyberNetwork/reserve-data/v3/common"
)

// isEthereumAddress is a validator.Func function that returns true if given field
// is a valid Ethereum address.
func isEthereumAddress(_ *validator.Validate, _ reflect.Value, _ reflect.Value,
	field reflect.Value, _ reflect.Type, _ reflect.Kind, _ string) bool {
	address := field.String()
	if len(address) != 0 && !common.IsHexAddress(address) {
		return false
	}
	return true
}

func validateCreateAsset(v *validator.Validate, sl *validator.StructLevel) {

	r := sl.CurrentStruct.Interface().(v3.CreateAssetEntry)
	if !r.IsQuote && v3.IsZeroAddress(r.Address) {
		sl.ReportError(reflect.ValueOf(r.Address), "address", "address", "address is required for non-quote")
	}
	if r.SetRate != v3.SetRateNotSet && r.PWI == nil {
		sl.ReportError(reflect.ValueOf(r.PWI), "PWI", "pwi", "pwi is required for setRate set")
	}
	if r.Rebalance && r.RebalanceQuadratic == nil {
		sl.ReportError(reflect.ValueOf(r.RebalanceQuadratic), "RebalanceQuadratic", "rebalance_quadratic", "rebalance_quadratic required with rebalance")
	}
	if r.Rebalance && r.Target == nil {
		sl.ReportError(reflect.ValueOf(r.Target), "Target", "target", "target required for rebalance")
	}
	if r.Transferable {
		for _, ex := range r.Exchanges {
			if v3.IsZeroAddress(ex.DepositAddress) {
				sl.ReportError(reflect.ValueOf(r.Transferable), "Transferable", "transferable", "transferable require deposit address")
				break
			}
		}
	}
}
func validateTradingPair(v *validator.Validate, sl *validator.StructLevel) {
	r := sl.CurrentStruct.Interface().(v3.TradingPair)
	if r.Base == 0 && r.Quote == 0 {
		sl.ReportError(reflect.ValueOf(r.Base), "base+quote", "base and quote", "base and quote rule")
	}
	if r.Base != 0 && r.Quote != 0 {
		sl.ReportError(reflect.ValueOf(r.Base), "base+quote", "base and quote", "base and quote rule")
	}
}
func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		var validators = []struct {
			name string
			fn   validator.Func
		}{
			{"isAddress", isEthereumAddress},
		}
		for _, val := range validators {
			if err := v.RegisterValidation(val.name, val.fn); err != nil {
				panic(err)
			}
		}
		v.RegisterStructValidation(validateCreateAsset, v3.CreateAssetEntry{})
		v.RegisterStructValidation(validateTradingPair, v3.TradingPair{})
	}
}
