package blockchain

import (
	"fmt"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli"

	"github.com/KyberNetwork/reserve-data/cmd/deployment"
	"github.com/KyberNetwork/reserve-data/v3/common"
)

const (
	wrapperAddressFlag = "wrapper-address"
	rateAddressFlag    = "rate-address"
	ethereumNodeFlag   = "ethereum-node"
	infuraEndpoint     = "https://mainnet.infura.io"
)

var (
	defaultWrapperAddress = map[deployment.Deployment]string{
		deployment.Development: "0x6172AFC8c00c46E0D07ce3AF203828198194620a",
		deployment.Ropsten:     "0x9de0a60F4A489e350cD8E3F249f4080858Af41d3",
		deployment.Production:  "0x6172AFC8c00c46E0D07ce3AF203828198194620a",
	}

	defaultRateAddress = map[deployment.Deployment]string{
		deployment.Development: "0x798AbDA6Cc246D0EDbA912092A2a3dBd3d11191B",
		deployment.Ropsten:     "0x535DE1F5a982c2a896da62790a42723A71c0c12B",
		deployment.Staging:     "0xe3E415a7a6c287a95DC68a01ff036828073fD2e6",
		deployment.Production:  "0x798AbDA6Cc246D0EDbA912092A2a3dBd3d11191B",
	}
)

// NewWrapperAddressFlag return wrapper address flag
func NewWrapperAddressFlag() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   wrapperAddressFlag,
			Usage:  "wrapper address",
			EnvVar: "WRAPPER_ADDRESS",
		},
		cli.StringFlag{
			Name:   rateAddressFlag,
			Usage:  "rate address",
			EnvVar: "RATE_ADDRESS",
		},
	}
}

// NewWrapperAddressFromContext return wrapper address and error if any
func NewWrapperAddressFromContext(c *cli.Context) (ethereum.Address, error) {
	wrapperAddress := ethereum.HexToAddress(c.String(wrapperAddressFlag))
	if common.IsZeroAddress(wrapperAddress) {
		dpl, err := deployment.NewDeploymentFromContext(c)
		if err != nil {
			return ethereum.Address{}, err
		}
		if address, exist := defaultWrapperAddress[dpl]; exist {
			wrapperAddress = ethereum.HexToAddress(address)
		} else {
			return ethereum.Address{}, fmt.Errorf("deployment does not have default wrapper address value: %s", dpl.String())
		}
	}
	return wrapperAddress, nil
}

// NewRateAddressFromContext return rate address and error if any
func NewRateAddressFromContext(c *cli.Context) (ethereum.Address, error) {
	rateAddress := ethereum.HexToAddress(c.String(rateAddressFlag))
	if common.IsZeroAddress(rateAddress) {
		dpl, err := deployment.NewDeploymentFromContext(c)
		if err != nil {
			return ethereum.Address{}, err
		}
		if address, exist := defaultRateAddress[dpl]; exist {
			rateAddress = ethereum.HexToAddress(address)
		} else {
			return ethereum.Address{}, fmt.Errorf("deployment does not have default rate address value: %s", dpl.String())
		}
	}
	return rateAddress, nil
}

// NewEthereumNodeFlags returns cli flag for ethereum node url input
func NewEthereumNodeFlags() cli.Flag {
	return cli.StringFlag{
		Name:   ethereumNodeFlag,
		Usage:  "Ethereum Node URL",
		EnvVar: "ETHEREUM_NODE",
		Value:  infuraEndpoint,
	}
}

// NewEthereumClientFromFlag returns Ethereum client from flag variable, or error if occurs
func NewEthereumClientFromFlag(c *cli.Context) (*ethclient.Client, error) {
	ethereumNodeURL := c.GlobalString(ethereumNodeFlag)
	return ethclient.Dial(ethereumNodeURL)
}
