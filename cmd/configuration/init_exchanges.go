package configuration

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"

	ethereum "github.com/ethereum/go-ethereum/common"

	"github.com/KyberNetwork/reserve-data/cmd/deployment"
	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/common/blockchain"
	"github.com/KyberNetwork/reserve-data/common/blockchain/nonce"
	"github.com/KyberNetwork/reserve-data/data/fetcher"
	"github.com/KyberNetwork/reserve-data/exchange"
	"github.com/KyberNetwork/reserve-data/exchange/binance"
	"github.com/KyberNetwork/reserve-data/exchange/huobi"
	"github.com/KyberNetwork/reserve-data/v3/storage"
)

type ExchangePool struct {
	Exchanges map[common.ExchangeID]interface{}
}

func updateDepositAddress(
	assetStorage storage.Interface,
	be exchange.BinanceInterface,
	he exchange.HuobiInterface) {
	assets, err := assetStorage.GetTransferableAssets()
	if err != nil {
		log.Printf("failed to get transferable assets err=%s", err.Error())
		return
	}
	for _, asset := range assets {
		for _, ae := range asset.Exchanges {
			switch ae.ExchangeID {
			case uint64(common.Binance):
				log.Printf("updating deposit address for asset id=%d exchange=%s symbol=%s",
					asset.ID,
					common.Binance.String(),
					ae.Symbol)
				depositAddress, err := be.GetDepositAddress(ae.Symbol)
				if err != nil {
					log.Printf("failed to get deposit address for asset id=%d exchange=%s symbol=%s err=%s",
						asset.ID,
						common.Binance.String(),
						ae.Symbol, err.Error())
					continue
				}
				err = assetStorage.UpdateDepositAddress(
					asset.ID,
					uint64(common.Binance),
					ethereum.HexToAddress(depositAddress.Address))
				if err != nil {
					log.Printf("assetStorage.UpdateDepositAddress err=%s", err.Error())
					continue
				}
			case uint64(common.Huobi):
				log.Printf("updating deposit address for asset id=%d exchange=%s symbol=%s",
					asset.ID,
					common.Huobi.String(),
					ae.Symbol)
				depositAddress, err := he.GetDepositAddress(ae.Symbol)
				if err != nil {
					log.Printf("failed to get deposit address for asset id=%d exchange=%s symbol=%s err=%s",
						asset.ID,
						common.Huobi.String(),
						ae.Symbol, err.Error())
					continue
				}
				err = assetStorage.UpdateDepositAddress(
					asset.ID,
					uint64(common.Huobi),
					ethereum.HexToAddress(depositAddress.Address))
				if err != nil {
					log.Printf("assetStorage.UpdateDepositAddress err=%s", err.Error())
					continue
				}
			}
		}
	}
}

func NewExchangePool(
	secretConfigFile string,
	blockchain *blockchain.BaseBlockchain,
	dpl deployment.Deployment,
	bi binance.Interface,
	hi huobi.Interface,
	enabledExchanges []common.ExchangeName,
	assetStorage storage.Interface,
) (*ExchangePool, error) {
	exchanges := map[common.ExchangeID]interface{}{}
	var (
		be *binance.Endpoint
		he *huobi.Endpoint
	)
	for _, exparam := range enabledExchanges {
		switch exparam {
		case common.StableExchange:
			stableEx, err := exchange.NewStableEx()
			if err != nil {
				return nil, fmt.Errorf("can not create exchange stable_exchange: (%s)", err.Error())
			}
			exchanges[stableEx.ID()] = stableEx
		case common.Binance:
			binanceSigner := binance.NewSignerFromFile(secretConfigFile)
			be = binance.NewBinanceEndpoint(binanceSigner, bi, dpl)
			binanceStorage, err := binance.NewBoltStorage(filepath.Join(common.CmdDirLocation(), "binance.db"))
			if err != nil {
				return nil, fmt.Errorf("can not create Binance storage: (%s)", err.Error())
			}
			bin, err := exchange.NewBinance(
				be,
				binanceStorage,
				assetStorage)
			if err != nil {
				return nil, fmt.Errorf("can not create exchange Binance: (%s)", err.Error())
			}
			exchanges[bin.ID()] = bin
		case common.Huobi:
			huobiSigner := huobi.NewSignerFromFile(secretConfigFile)
			he = huobi.NewHuobiEndpoint(huobiSigner, hi)
			huobiStorage, err := huobi.NewBoltStorage(filepath.Join(common.CmdDirLocation(), "huobi.db"))
			if err != nil {
				return nil, fmt.Errorf("can not create Huobi storage: (%s)", err.Error())
			}
			intermediatorSigner, err := HuobiIntermediatorSignerFromFile(secretConfigFile)
			if err != nil {
				log.Printf("failed to get itermediator signer from file err=%s", err.Error())
				return nil, err
			}
			intermediatorNonce := nonce.NewTimeWindow(intermediatorSigner.GetAddress(), 10000)
			hb, err := exchange.NewHuobi(
				he,
				blockchain,
				intermediatorSigner,
				intermediatorNonce,
				huobiStorage,
				assetStorage,
			)
			if err != nil {
				return nil, fmt.Errorf("can not create exchange Huobi: (%s)", err.Error())
			}
			exchanges[hb.ID()] = hb
		}
	}

	go updateDepositAddress(assetStorage, be, he)
	return &ExchangePool{exchanges}, nil
}

func (ep *ExchangePool) FetcherExchanges() ([]fetcher.Exchange, error) {
	result := []fetcher.Exchange{}
	for _, ex := range ep.Exchanges {
		fcEx, ok := ex.(fetcher.Exchange)
		if !ok {
			return result, errors.New("ExchangePool cannot be asserted  to fetcher exchange")
		}
		result = append(result, fcEx)
	}
	return result, nil
}

func (ep *ExchangePool) CoreExchanges() ([]common.Exchange, error) {
	result := []common.Exchange{}
	for _, ex := range ep.Exchanges {
		cmEx, ok := ex.(common.Exchange)
		if !ok {
			return result, errors.New("ExchangePool cannot be asserted to core Exchange")
		}
		result = append(result, cmEx)
	}
	return result, nil
}
