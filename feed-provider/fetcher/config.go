package fetcher

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/urfave/cli"
)

const (
	configLocationFlag = "config-file"
	defaultLocation    = "config.json"
)

// Config base config for fetcher
type Config struct {
	URL    string  `json:"url"`
	Amount float64 `json:"amount"`
}

// Configs config for init setup
type Configs struct {
	CoinbaseETHDAI10000 Config `json:"CoinbaseETHDAI10000"`
	KrakenETHDAI10000   Config `json:"KrakenETHDAI10000"`
	CoinbaseETHBTC3     Config `json:"CoinbaseETHBTC3"`
	BinanceETHBTC3      Config `json:"BinanceETHBTC3"`
}

// NewConfigCliFlags flag to config location
func NewConfigCliFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   configLocationFlag,
			Usage:  "location to config file",
			EnvVar: "FEED_PROVIDER_CONFIG",
			Value:  defaultLocation,
		},
	}
}

// NewConfigFromCli load config from file and create new instance
func NewConfigFromCli(c *cli.Context) (*Configs, error) {
	fileLocation := c.String(configLocationFlag)
	jsonFile, err := os.Open(fileLocation)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var configs Configs
	err = json.Unmarshal(byteValue, &configs)
	return &configs, err
}
