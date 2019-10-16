package blockchain

import (
	"context"
	"errors"
	"math/big"
	"time"

	ether "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

type ContractCaller struct {
	clients []*ethclient.Client
	urls    []string
	l       *zap.SugaredLogger
}

func NewContractCaller(clients []*ethclient.Client, urls []string) *ContractCaller {
	return &ContractCaller{
		clients: clients,
		urls:    urls,
		l:       zap.S(),
	}
}

func (c ContractCaller) CallContract(msg ether.CallMsg, blockNo *big.Int, timeOut time.Duration) ([]byte, error) {
	for i, client := range c.clients {
		url := c.urls[i]

		output, err := func() ([]byte, error) {
			ctx, cancel := context.WithTimeout(context.Background(), timeOut)
			defer cancel()
			return client.CallContract(ctx, msg, blockNo)
		}()
		if err != nil {
			c.l.Warnf("FALLBACK: Ether client %s done, getting err %v, trying next one...", url, err)
			continue
		}
		return output, nil
	}
	return nil, errors.New("failed to call contract, all clients failed")
}
