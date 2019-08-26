package blockchain

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/KyberNetwork/reserve-data/v3/blockchain/contracts"
)

// Blockchain to interact with blockchain
type Blockchain struct {
	Wrapper     *contracts.Wrapper
	RateAddress common.Address
}

// NewBlockchain return new blockchain instance
func NewBlockchain(wrapperAddress, rateAddress common.Address, ethClient *ethclient.Client) (*Blockchain, error) {
	wrapper, err := contracts.NewWrapper(wrapperAddress, ethClient)
	if err != nil {
		return nil, err
	}
	return &Blockchain{
		Wrapper:     wrapper,
		RateAddress: rateAddress,
	}, nil
}

// CheckTokenIndices check if a token is listed on reserve or not
func (bc *Blockchain) CheckTokenIndices(tokenAddr common.Address) error {
	tokenAddrs := []common.Address{}
	tokenAddrs = append(tokenAddrs, tokenAddr)
	opts := &bind.CallOpts{BlockNumber: big.NewInt(0)}
	_, _, err := bc.Wrapper.GetTokenIndicies(
		opts,
		bc.RateAddress,
		tokenAddrs,
	)
	if err != nil {
		return err
	}
	return nil
}
