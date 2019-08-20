package configuration

import (
	"fmt"

	"github.com/urfave/cli"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/exchange/binance"
	"github.com/KyberNetwork/reserve-data/exchange/huobi"
)

const (
	exchangesFlag        = "exchanges"
	binanceAPIKeyFlag    = "binance-api-key"
	binanceSecretKeyFlag = "binance-secret-key"
	huobiAPIKeyFlag      = "huobi-api-key"
	huobiSecretKeyFlag   = "huobi-secret-key"
)

// NewBinanceKeysFlag add flag for binance keys
func NewBinanceKeysFlag() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   binanceAPIKeyFlag,
			Usage:  "api key for binance",
			EnvVar: "BINANCE_API_KEY",
		},
		cli.StringFlag{
			Name:   binanceSecretKeyFlag,
			Usage:  "secret key for binance",
			EnvVar: "BINANCE_SECRET_KEY",
		},
	}
}

// NewBinanceSignerFromContext return binance signer
func NewBinanceSignerFromContext(c *cli.Context) binance.Signer {
	apiKey := c.String(binanceAPIKeyFlag)
	secretKey := c.String(binanceSecretKeyFlag)
	return binance.NewSigner(apiKey, secretKey)
}

// NewHuobiKeysFlag add flag for huobi keys
func NewHuobiKeysFlag() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   huobiAPIKeyFlag,
			Usage:  "api key for huobi",
			EnvVar: "HUOBI_API_KEY",
		},
		cli.StringFlag{
			Name:   huobiSecretKeyFlag,
			Usage:  "secret key for huobi",
			EnvVar: "HUOBI_SECRET_KEY",
		},
	}
}

// NewHuobiSignerFromContext return huobi signer from context
func NewHuobiSignerFromContext(c *cli.Context) huobi.Signer {
	apiKey := c.String(huobiAPIKeyFlag)
	secretKey := c.String(huobiSecretKeyFlag)
	return huobi.NewSigner(apiKey, secretKey)
}

// NewExchangeCliFlag creates new cli flag for enable/disable exchanges.
func NewExchangeCliFlag() cli.Flag {
	return cli.StringSliceFlag{
		Name:   exchangesFlag,
		Usage:  "Enable an exchange for fetching order books, rebalancing purpose. By default all exchanges are disabled.",
		EnvVar: "KYBER_EXCHANGES",
	}
}

// NewExchangesFromContext returns configured exchanges from cli context.
func NewExchangesFromContext(c *cli.Context) ([]common.ExchangeID, error) {
	var exchanges []common.ExchangeID

	for _, exchangeName := range c.GlobalStringSlice(exchangesFlag) {
		exchange, ok := common.ValidExchangeNames[exchangeName]
		if !ok {
			return nil, fmt.Errorf("invalid exchange %v", exchangeName)
		}
		exchanges = append(exchanges, exchange)
	}

	return exchanges, nil
}
