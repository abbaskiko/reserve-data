package blockchain

import (
	"log"
	"testing"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KyberNetwork/reserve-data/common"
	baseblockchain "github.com/KyberNetwork/reserve-data/common/blockchain"
)

const (
	rpcEndpoint = "https://mainnet.infura.io/v3/3243b06bf0334cff8e468bf90ce6e08c"
)

// this is an external test
func TestGeneratedGetListedTokens(t *testing.T) {
	// t.Skip() // skip as it is an external test
	log.Println("test generated get listed tokens")

	rpcClient, err := rpc.Dial(rpcEndpoint)
	require.NoError(t, err)

	ethClient := ethclient.NewClient(rpcClient)
	commonEthClient := common.EthClient{
		Client: ethClient,
		URL:    rpcEndpoint,
	}
	contractCaller := baseblockchain.NewContractCaller([]*common.EthClient{&commonEthClient})

	baseBlockchain := baseblockchain.NewBaseBlockchain(
		rpcClient,      // rpc client
		ethClient,      // eth client
		nil,            // map[string]*Operator
		nil,            // broadcaster
		contractCaller, // contract caller
	)
	contracts := common.ContractAddressConfiguration{
		Reserve: ethereum.Address{},
		Wrapper: ethereum.Address{},
		Pricing: ethereum.HexToAddress("0x798AbDA6Cc246D0EDbA912092A2a3dBd3d11191B"),
		Proxy:   ethereum.Address{},
	}

	blockchain, err := NewBlockchain(baseBlockchain, &contracts, nil, common.GasConfig{})
	require.NoError(t, err)

	opts := blockchain.GetCallOpts(0)
	listedToken, err := blockchain.GeneratedGetListedTokens(opts)
	assert.NoError(t, err)
	assert.NotEmpty(t, listedToken)
	log.Printf("listed tokens: %d", len(listedToken))
	for _, token := range listedToken {
		log.Printf("token: %s", token.Hex())
	}
}
