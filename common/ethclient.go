package common

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// EthClient just a wrap for ethclient.Client with url connect to
type EthClient struct {
	*ethclient.Client
	URL       string
	RPCClient *rpc.Client
}

// NewEthClient create a new ethclient
func NewEthClient(url string) (*EthClient, error) {
	r, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	c := ethclient.NewClient(r)
	return &EthClient{
		Client:    c,
		URL:       url,
		RPCClient: r,
	}, nil
}
