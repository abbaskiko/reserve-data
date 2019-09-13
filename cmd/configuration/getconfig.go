package configuration

import (
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/urfave/cli"

	"github.com/KyberNetwork/reserve-data/cmd/deployment"
	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/common/archive"
	"github.com/KyberNetwork/reserve-data/common/blockchain"
	"github.com/KyberNetwork/reserve-data/exchange/binance"
	"github.com/KyberNetwork/reserve-data/exchange/huobi"
	"github.com/KyberNetwork/reserve-data/reservesetting/storage"
	"github.com/KyberNetwork/reserve-data/world"
)

const (
	byzantiumChainType = "byzantium"
	homesteadChainType = "homestead"
)

// GetChainType return chain type
func GetChainType(dpl deployment.Deployment) string {
	switch dpl {
	case deployment.Production:
		return byzantiumChainType
	case deployment.Development:
		return homesteadChainType
	case deployment.Kovan:
		return homesteadChainType
	case deployment.Staging:
		return byzantiumChainType
	case deployment.Simulation, deployment.Analytic:
		return homesteadChainType
	case deployment.Ropsten:
		return byzantiumChainType
	default:
		return homesteadChainType
	}
}

// GetConfig return config for core
func GetConfig(
	cliCtx *cli.Context,
	dpl deployment.Deployment,
	nodeConf *EthereumNodeConfiguration,
	bi binance.Interface,
	hi huobi.Interface,
	contractAddressConf *common.ContractAddressConfiguration,
	dataFile string,
	secretConfigFile string,
	settingStorage storage.Interface,
) (*Config, error) {
	theWorld, err := world.NewTheWorld(dpl, secretConfigFile)
	if err != nil {
		log.Printf("Can't init the world (which is used to get global data), err %s", err.Error())
		return nil, err
	}

	chainType := GetChainType(dpl)

	//set client & endpoint
	client, err := rpc.Dial(nodeConf.Main)
	if err != nil {
		return nil, err
	}

	mainClient := ethclient.NewClient(client)
	bkClients := map[string]*ethclient.Client{}

	var callClients []*ethclient.Client
	for _, ep := range nodeConf.Backup {
		var bkClient *ethclient.Client
		bkClient, err = ethclient.Dial(ep)
		if err != nil {
			log.Printf("Cannot connect to %s, err %s. Ignore it.", ep, err)
		} else {
			bkClients[ep] = bkClient
			callClients = append(callClients, bkClient)
		}
	}

	bc := blockchain.NewBaseBlockchain(
		client, mainClient, map[string]*blockchain.Operator{},
		blockchain.NewBroadcaster(bkClients),
		chainType,
		blockchain.NewContractCaller(callClients, nodeConf.Backup),
	)

	awsConf, err := archive.GetAWSconfigFromFile(secretConfigFile)
	if err != nil {
		log.Printf("failed to load AWS config from file %s", secretConfigFile)
		return nil, err
	}
	s3archive := archive.NewS3Archive(awsConf)
	config := &Config{
		Blockchain:              bc,
		EthereumEndpoint:        nodeConf.Main,
		BackupEthereumEndpoints: nodeConf.Backup,
		Archive:                 s3archive,
		World:                   theWorld,
		ContractAddresses:       contractAddressConf,
		SettingStorage:          settingStorage,
	}

	log.Printf("configured endpoint: %s, backup: %v", config.EthereumEndpoint, config.BackupEthereumEndpoints)
	if err = config.AddCoreConfig(cliCtx, secretConfigFile, dpl, bi, hi, contractAddressConf, dataFile, settingStorage); err != nil {
		return nil, err
	}
	return config, nil
}
