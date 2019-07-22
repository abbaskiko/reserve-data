package common

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// IsZeroAddress return if hash is zero.
func IsZeroAddress(addr common.Address) bool {
	return addr.Hash().Big().Cmp(big.NewInt(0)) == 0
}

// FloatPointer is helper use in optional parameter
func FloatPointer(f float64) *float64 {
	return &f
}

// StringPointer is helper use in optional parameter
func StringPointer(s string) *string {
	return &s
}

// BoolPointer is helper use in optional parameter
func BoolPointer(b bool) *bool {
	return &b
}
