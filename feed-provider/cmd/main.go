package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-data/feed-provider/collector"
	"github.com/KyberNetwork/reserve-data/feed-provider/fetcher"
	"github.com/KyberNetwork/reserve-data/feed-provider/fetcher/binance"
	"github.com/KyberNetwork/reserve-data/feed-provider/fetcher/coinbase"
	"github.com/KyberNetwork/reserve-data/feed-provider/fetcher/kraken"
	"github.com/KyberNetwork/reserve-data/feed-provider/httpserver"
	"github.com/KyberNetwork/reserve-data/feed-provider/storage"
	"github.com/KyberNetwork/reserve-data/lib/app"
	"github.com/KyberNetwork/reserve-data/lib/httputil"
)

func main() {
	app := cli.NewApp()
	app.Name = "Feed provider"
	app.Usage = "Calculate and convert feed data to keep the same format"
	app.Version = "0.0.1"
	app.Action = run

	app.Flags = NewCliFlags()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	l, flusher, err := app.NewSugaredLogger(c)
	if err != nil {
		panic(err)
	}
	defer func() {
		flusher()
	}()
	zap.ReplaceGlobals(l.Desugar())

	config, err := fetcher.NewConfigFromCli(c)
	if err != nil {
		return err
	}

	coinbaseETHDAI10000Fetcher, err := coinbase.NewFetcher(config.CoinbaseETHDAI10000, l)
	if err != nil {
		return err
	}
	krakenETHDAI10000Fetcher, err := kraken.NewFetcher(config.KrakenETHDAI10000, l)
	if err != nil {
		return err
	}
	coinbaseETHBTC3Fetcher, err := coinbase.NewFetcher(config.CoinbaseETHBTC3, l)
	if err != nil {
		return err
	}
	binanceETHBTC3Fetcher, err := binance.NewFetcher(config.BinanceETHBTC3, l)
	if err != nil {
		return err
	}

	fetchers := map[string]fetcher.Fetcher{
		"CoinbaseETHDAI10000": coinbaseETHDAI10000Fetcher,
		"KrakenETHDAI10000":   krakenETHDAI10000Fetcher,
		"CoinbaseETHBTC3":     coinbaseETHBTC3Fetcher,
		"BinanceETHBTC3":      binanceETHBTC3Fetcher,
	}

	s := storage.NewRAMStorage()
	host := httputil.NewHTTPAddressFromContext(c)
	server, err := httpserver.NewHTTPServer(l, s, host)
	if err != nil {
		return err
	}
	collector := collector.NewCollectorFromCli(c, l, s, fetchers)
	go collector.Run()
	return server.Run()
}
