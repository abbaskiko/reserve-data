package configuration

import (
	"fmt"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli"

	"github.com/KyberNetwork/reserve-data/cmd/deployment"
	"github.com/KyberNetwork/reserve-data/common"
)

const (
	reserveContractFlag         = "reserve-contract"
	wrapperContractFlag         = "wrapper-contract"
	pricingContractFlag         = "pricing-contract"
	internalNetworkContractFlag = "internal-network-contract"
	networkContractFlag         = "network-contract"
)

//defaultAddressConfiguration store token configs according to env mode.
var defaultAddressConfiguration = map[deployment.Deployment]common.ContractAddressConfiguration{
	deployment.Development: {
		Reserve:         ethereum.HexToAddress("0x63825c174ab367968EC60f061753D3bbD36A0D8F"),
		Wrapper:         ethereum.HexToAddress("0x6172AFC8c00c46E0D07ce3AF203828198194620a"),
		Pricing:         ethereum.HexToAddress("0x798AbDA6Cc246D0EDbA912092A2a3dBd3d11191B"),
		InternalNetwork: ethereum.HexToAddress("0x91a502C678605fbCe581eae053319747482276b9"),
		Network:         ethereum.HexToAddress("0x818E6FECD516Ecc3849DAf6845e3EC868087B755"),
	},
	deployment.Staging: {
		Reserve:         ethereum.HexToAddress("0x2C5a182d280EeB5824377B98CD74871f78d6b8BC"),
		Wrapper:         ethereum.HexToAddress("0x6172AFC8c00c46E0D07ce3AF203828198194620a"),
		Pricing:         ethereum.HexToAddress("0xe3E415a7a6c287a95DC68a01ff036828073fD2e6"),
		InternalNetwork: ethereum.HexToAddress("0x706aBcE058DB29eB36578c463cf295F180a1Fe9C"),
		Network:         ethereum.HexToAddress("0xC14f34233071543E979F6A79AA272b0AB1B4947D"),
	},
	deployment.Production: {
		Reserve:         ethereum.HexToAddress("0x63825c174ab367968EC60f061753D3bbD36A0D8F"),
		Wrapper:         ethereum.HexToAddress("0x6172AFC8c00c46E0D07ce3AF203828198194620a"),
		Pricing:         ethereum.HexToAddress("0x798AbDA6Cc246D0EDbA912092A2a3dBd3d11191B"),
		InternalNetwork: ethereum.HexToAddress("0x65bF64Ff5f51272f729BDcD7AcFB00677ced86Cd"),
		Network:         ethereum.HexToAddress("0x818E6FECD516Ecc3849DAf6845e3EC868087B755"),
	},
	deployment.Ropsten: {
		Reserve: ethereum.HexToAddress("0x0FC1CF3e7DD049F7B42e6823164A64F76fC06Be0"),
		Wrapper: ethereum.HexToAddress("0x9de0a60F4A489e350cD8E3F249f4080858Af41d3"),
		Pricing: ethereum.HexToAddress("0x535DE1F5a982c2a896da62790a42723A71c0c12B"),
		Network: ethereum.HexToAddress("0x0a56d8a49E71da8d7F9C65F95063dB48A3C9560B"),
	},
}

// EthereumNodeConfiguration contains all Ethereum nodes to interactive with Ethereum network.
type EthereumNodeConfiguration struct {
	Main   string
	Backup []string
}

// NewEthereumNodeConfiguration returns a new Ethereum node configuration.
func NewEthereumNodeConfiguration(main string, backup []string) *EthereumNodeConfiguration {
	return &EthereumNodeConfiguration{Main: main, Backup: backup}
}

// NewContractAddressCliFlags returns new cli flags for address contract configuration.
func NewContractAddressCliFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   reserveContractFlag,
			Usage:  "Reserve contract address",
			EnvVar: "RESERVE_CONTRACT",
		},
		cli.StringFlag{
			Name:   wrapperContractFlag,
			Usage:  "Wrapper contract address",
			EnvVar: "WRAPPER_CONTRACT",
		},
		cli.StringFlag{
			Name:   pricingContractFlag,
			Usage:  "Pricing contract address",
			EnvVar: "PRICING_CONTRACT",
		},
		cli.StringFlag{
			Name:   internalNetworkContractFlag,
			Usage:  "internal network contract address",
			EnvVar: "INTERNAL_NETWORK_CONTRACT",
		},
		cli.StringFlag{
			Name:   networkContractFlag,
			Usage:  "network contract address",
			EnvVar: "NETWORK_CONTRACT",
		},
	}
}

// NewContractAddressConfigurationFromContext returns the contract address configuration from cli flags.
func NewContractAddressConfigurationFromContext(c *cli.Context) (*common.ContractAddressConfiguration, error) {
	var (
		conf                       = &common.ContractAddressConfiguration{}
		reserveContractStr         = c.GlobalString(reserveContractFlag)
		wrapperContractStr         = c.GlobalString(wrapperContractFlag)
		pricingContractStr         = c.GlobalString(pricingContractFlag)
		internalNetworkContractStr = c.GlobalString(internalNetworkContractFlag)
		networkContractStr         = c.GlobalString(networkContractFlag)
	)

	if len(reserveContractStr) != 0 && len(wrapperContractStr) != 0 && len(pricingContractStr) != 0 {
		if !ethereum.IsHexAddress(reserveContractStr) {
			return nil, fmt.Errorf("invalid reserve contract address %s", reserveContractStr)
		}
		if !ethereum.IsHexAddress(wrapperContractStr) {
			return nil, fmt.Errorf("invalid wrapper contract address %s", wrapperContractStr)
		}
		if !ethereum.IsHexAddress(pricingContractStr) {
			return nil, fmt.Errorf("invalid pricing contract address %s", pricingContractStr)
		}
		conf.Reserve = ethereum.HexToAddress(reserveContractStr)
		conf.Wrapper = ethereum.HexToAddress(wrapperContractStr)
		conf.Pricing = ethereum.HexToAddress(pricingContractStr)
		conf.InternalNetwork = ethereum.HexToAddress(internalNetworkContractStr)
		conf.Network = ethereum.HexToAddress(networkContractStr)
		return conf, nil
	}

	dpl, err := deployment.NewDeploymentFromContext(c)
	if err != nil {
		return nil, err
	}

	defaultConf, ok := defaultAddressConfiguration[dpl]
	if !ok {
		return nil, fmt.Errorf("no default contract addresses configuration for deployment %s", dpl.String())
	}

	if len(reserveContractStr) == 0 {
		conf.Reserve = defaultConf.Reserve
	} else {
		if !ethereum.IsHexAddress(reserveContractStr) {
			return nil, fmt.Errorf("invalid reserve contract address %s", reserveContractStr)
		}
		conf.Reserve = ethereum.HexToAddress(reserveContractStr)
	}

	if len(wrapperContractStr) == 0 {
		conf.Wrapper = defaultConf.Wrapper
	} else {
		if !ethereum.IsHexAddress(wrapperContractStr) {
			return nil, fmt.Errorf("invalid wrapper contract address %s", wrapperContractStr)
		}
		conf.Wrapper = ethereum.HexToAddress(wrapperContractStr)
	}

	if len(pricingContractStr) == 0 {
		conf.Pricing = defaultConf.Pricing
	} else {
		if !ethereum.IsHexAddress(pricingContractStr) {
			return nil, fmt.Errorf("invalid pricing contract address %s", pricingContractStr)
		}
		conf.Pricing = ethereum.HexToAddress(pricingContractStr)
	}

	if len(internalNetworkContractStr) == 0 {
		conf.InternalNetwork = defaultConf.InternalNetwork
	} else {
		if !ethereum.IsHexAddress(internalNetworkContractStr) {
			return nil, fmt.Errorf("invalid pricing contract address %s", internalNetworkContractStr)
		}
		conf.InternalNetwork = ethereum.HexToAddress(internalNetworkContractStr)
	}

	if len(networkContractStr) == 0 {
		conf.Network = defaultConf.Network
	} else {
		if !ethereum.IsHexAddress(networkContractStr) {
			return nil, fmt.Errorf("invalid pricing contract address %s", networkContractStr)
		}
		conf.Network = ethereum.HexToAddress(networkContractStr)
	}

	return conf, nil
}
