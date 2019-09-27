package common

import (
	"github.com/ethereum/go-ethereum/ethclient"
)

// EthClient just a wrap for ethclient.Client with url connect to
type EthClient struct {
	*ethclient.Client
	URL string
}
