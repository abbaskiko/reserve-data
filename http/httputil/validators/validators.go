package validators

import (
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v8"
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

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		var validators = []struct {
			name string
			fn   validator.Func
		}{
			{"isEthereumAddress", isEthereumAddress},
		}
		for _, val := range validators {
			if err := v.RegisterValidation(val.name, val.fn); err != nil {
				panic(err)
			}
		}
	}
}
