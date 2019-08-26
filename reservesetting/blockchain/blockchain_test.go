package blockchain

import (
	"testing"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckTokenIndices(t *testing.T) {
	//? this is an which call to external resourc
	//? skip after run in local
	t.Skip()

	wrapperAddress := ethereum.HexToAddress("0x6172AFC8c00c46E0D07ce3AF203828198194620a")
	rateAddress := ethereum.HexToAddress("0x798AbDA6Cc246D0EDbA912092A2a3dBd3d11191B")

	ethClient, err := ethclient.Dial("https://mainnet.infura.io/v3/fd330878eff84d48b97e3023c996dff6")
	require.NoError(t, err)

	blockchain, err := NewBlockchain(wrapperAddress, rateAddress, ethClient)
	require.NoError(t, err)

	notIndexAddress := ethereum.HexToAddress("0x3c513823517Ed8BFa9e93Bf8840840BeDBF52821")

	err = blockchain.CheckTokenIndices(notIndexAddress)
	assert.Error(t, err)

	indexAddress := ethereum.HexToAddress("0xdd974d5c2e2928dea5f71b9825b8b646686bd200")
	err = blockchain.CheckTokenIndices(indexAddress)
	assert.NoError(t, err)
}
