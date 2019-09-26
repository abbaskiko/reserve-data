package blockchain

import (
	"context"
	"errors"
	"log"
	"math/big"
	"time"

	ether "github.com/ethereum/go-ethereum"

	"github.com/KyberNetwork/reserve-data/common"
)

type ContractCaller struct {
	clients []*common.EthClient
}

func NewContractCaller(clients []*common.EthClient) *ContractCaller {
	return &ContractCaller{
		clients: clients,
	}
}

func (c ContractCaller) CallContract(msg ether.CallMsg, blockNo *big.Int, timeOut time.Duration) ([]byte, error) {
	for _, client := range c.clients {

		output, err := func() ([]byte, error) {
			ctx, cancel := context.WithTimeout(context.Background(), timeOut)
			defer cancel()
			return client.CallContract(ctx, msg, blockNo)
		}()
		if err != nil {
			log.Printf("FALLBACK: Ether client %s done, getting err %v, trying next one...", client.Url, err)
			continue
		}
		return output, nil
	}
	return nil, errors.New("failed to call contract, all clients failed")
}
