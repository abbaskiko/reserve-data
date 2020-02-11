package blockchain

import (
	"bytes"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ethereum "github.com/ethereum/go-ethereum/common"
)

type Contract struct {
	Address ethereum.Address
	ABI     abi.ABI
}

func NewContract(address ethereum.Address, abiData string) *Contract {
	parsed, err := abi.JSON(bytes.NewBufferString(abiData))
	if err != nil {
		panic(err)
	}
	return &Contract{address, parsed}
}
