package common

import (
	ethereum "github.com/ethereum/go-ethereum/common"
)

// ContractAddressConfiguration contains the smart contract addresses.
type ContractAddressConfiguration struct {
	Reserve ethereum.Address
	// Wrapper is not officially shown but you can find it here:
	// https://github.com/KyberNetwork/smart-contracts/blob/master/contracts/mock/Wrapper.sol
	Wrapper ethereum.Address
	// Pricing is ConversionRates contract address
	Pricing ethereum.Address
	// Proxy contract address
	Proxy ethereum.Address
	// RateQueryHelper contract address
	RateQueryHelper ethereum.Address
}
