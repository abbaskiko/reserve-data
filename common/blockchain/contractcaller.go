package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"time"

	ether "github.com/ethereum/go-ethereum"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-data/common"
)

type ContractCaller struct {
	clients []*common.EthClient
	l       *zap.SugaredLogger
}

func NewContractCaller(clients []*common.EthClient) *ContractCaller {
	return &ContractCaller{
		clients: clients,
		l:       zap.S(),
	}
}

func (c ContractCaller) CallContract(msg ether.CallMsg, blockNo *big.Int, timeOut time.Duration) ([]byte, error) {
	type errInfo struct {
		URL string
		Err error
	}
	var errs []errInfo
	for _, client := range c.clients {
		output, err := func() ([]byte, error) {
			ctx, cancel := context.WithTimeout(context.Background(), timeOut)
			defer cancel()
			return client.CallContract(ctx, msg, blockNo)
		}()
		if err != nil {
			c.l.Infof("FALLBACK: Ether client %s done, getting err %v, trying next one...", client.URL, err)
			errs = append(errs, errInfo{
				URL: client.URL,
				Err: err,
			})
			continue
		}
		return output, nil
	}
	return nil, fmt.Errorf("failed to call contract, all clients failed: %v", errs)
}
