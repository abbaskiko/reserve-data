package common

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

//IsZeroAddress return if hash is zero.
func IsZeroAddress(addr common.Address) bool {
	return addr.Hash().Big().Cmp(big.NewInt(0)) == 0
}
