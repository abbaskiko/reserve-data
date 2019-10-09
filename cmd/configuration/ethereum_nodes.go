package configuration

import (
	"fmt"

	"github.com/urfave/cli"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-data/cmd/deployment"
)

const (
	ethereumMainNodeFlag   = "ethereum-main-node"
	ethereumBackupNodeFlag = "ethereum-backup-nodes"

	infuraProjectID = "/v3/59d9e06a1abe487e8e74664c06b337f9"

	// AlchemyapiMainnetEndpoint exported for reserve setting
	AlchemyapiMainnetEndpoint = "https://eth-mainnet.alchemyapi.io/jsonrpc/V1GjKybGLx6rzSu517KSWpSrTSIIXmV7"

	// InfuraMainnetEndpoint exported to use in reserve setting
	InfuraMainnetEndpoint = "https://mainnet.infura.io" + infuraProjectID
	// InfuraKovanEndpoint exported to use in reserve setting
	InfuraKovanEndpoint = "https://kovan.infura.io" + infuraProjectID
	// InfuraRopstenEndpoint exported to use in reserve setting
	InfuraRopstenEndpoint = "https://ropsten.infura.io" + infuraProjectID

	myEtherAPIMainnetEndpoint = "https://api.myetherwallet.com/eth"
	myEtherAPIRopstenEndpoint = "https://api.myetherwallet.com/rop"
	semidNodeKyberEndpoint    = "https://semi-node.kyber.network"
	myCryptoAPIEndpoint       = "https://api.mycryptoapi.com/eth"
	mewGivethAPIEndpoint      = "https://mew.giveth.io/"

	localDevChainEndpoint = "http://blockchain:8545"
)

var defaultEthereumNodes = map[deployment.Deployment]*EthereumNodeConfiguration{
	deployment.Development: NewEthereumNodeConfiguration(
		InfuraMainnetEndpoint,
		[]string{
			semidNodeKyberEndpoint,
			InfuraMainnetEndpoint,
			myCryptoAPIEndpoint,
			myEtherAPIMainnetEndpoint,
		},
	),
	deployment.Kovan: NewEthereumNodeConfiguration(
		InfuraKovanEndpoint,
		[]string{},
	),
	deployment.Production: NewEthereumNodeConfiguration(
		AlchemyapiMainnetEndpoint,
		[]string{
			AlchemyapiMainnetEndpoint,
			InfuraMainnetEndpoint,
			semidNodeKyberEndpoint,
			myCryptoAPIEndpoint,
			myEtherAPIMainnetEndpoint,
			mewGivethAPIEndpoint,
		},
	),
	deployment.Staging: NewEthereumNodeConfiguration(
		AlchemyapiMainnetEndpoint,
		[]string{
			AlchemyapiMainnetEndpoint,
			InfuraMainnetEndpoint,
			semidNodeKyberEndpoint,
			myCryptoAPIEndpoint,
			myEtherAPIMainnetEndpoint,
			mewGivethAPIEndpoint,
		},
	),
	deployment.Simulation: NewEthereumNodeConfiguration(
		localDevChainEndpoint,
		[]string{
			localDevChainEndpoint,
		},
	),
	deployment.Ropsten: NewEthereumNodeConfiguration(
		InfuraRopstenEndpoint,
		[]string{
			myEtherAPIRopstenEndpoint,
		},
	),
	deployment.Analytic: NewEthereumNodeConfiguration(
		localDevChainEndpoint,
		[]string{
			localDevChainEndpoint,
		},
	),
}

// NewEthereumNodesCliFlags returns new cli flag for config ethereum nodes.
func NewEthereumNodesCliFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   ethereumMainNodeFlag,
			Usage:  "Ethereum main node to use to interact with Ethereum network",
			EnvVar: "ETHEREUM_MAIN_NODE",
		},
		cli.StringSliceFlag{
			Name:   ethereumBackupNodeFlag,
			Usage:  "Ethereum nodes to use in broadcast transaction in case the main node is unreachable",
			EnvVar: "ETHEREUM_BACKUP_NODES",
		},
	}
}

// NewEthereumNodeConfigurationFromContext returns the configured ethereum node from cli context.
func NewEthereumNodeConfigurationFromContext(c *cli.Context, l *zap.SugaredLogger) (*EthereumNodeConfiguration, error) {
	var (
		conf        = &EthereumNodeConfiguration{}
		mainNode    = c.GlobalString(ethereumMainNodeFlag)
		backupNodes = c.StringSlice(ethereumBackupNodeFlag)
	)

	if len(mainNode) != 0 && len(backupNodes) != 0 {
		l.Infow("using provided Ethereum", "main_node", mainNode, "backup_node", backupNodes)
		return NewEthereumNodeConfiguration(mainNode, backupNodes), nil
	}

	dpl, err := deployment.NewDeploymentFromContext(c)
	if err != nil {
		return nil, err
	}

	defaultConf, ok := defaultEthereumNodes[dpl]
	if !ok {
		return nil, fmt.Errorf("no default ethereum node configuration for deployment %s", dpl.String())
	}

	if len(backupNodes) != 0 {
		l.Infow("using provided Ethereum backup nodes", "nodes", backupNodes)
		conf.Backup = backupNodes
	} else {
		conf.Backup = defaultConf.Backup
	}

	if len(mainNode) != 0 {
		l.Infow("using provided Ethereum main node", "main_node", mainNode)
		conf.Main = mainNode
		// transaction broadcasting only use backup nodes, prepend the provided node
		// to make sure it has the highest priority
		conf.Backup = append([]string{mainNode}, conf.Backup...)
	} else {
		conf.Main = defaultConf.Main
	}

	return conf, nil
}
