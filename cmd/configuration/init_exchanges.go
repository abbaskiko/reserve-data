package configuration

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"

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

func NewExchangePool(
	secretConfigFile string,
	blockchain *blockchain.BaseBlockchain,
	dpl deployment.Deployment,
	bi binance.Interface,
	hi huobi.Interface,
	enabledExchanges []common.ExchangeName,
	sr storage.SettingReader,
) (*ExchangePool, error) {
	exchanges := map[common.ExchangeID]interface{}{}
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
			endpoint := binance.NewBinanceEndpoint(binanceSigner, bi, dpl)
			binanceStorage, err := binance.NewBoltStorage(filepath.Join(common.CmdDirLocation(), "binance.db"))
			if err != nil {
				return nil, fmt.Errorf("can not create Binance storage: (%s)", err.Error())
			}
			bin, err := exchange.NewBinance(
				endpoint,
				binanceStorage,
				sr)
			if err != nil {
				return nil, fmt.Errorf("can not create exchange Binance: (%s)", err.Error())
			}
			exchanges[bin.ID()] = bin
		case common.Huobi:
			huobiSigner := huobi.NewSignerFromFile(secretConfigFile)
			endpoint := huobi.NewHuobiEndpoint(huobiSigner, hi)
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
				endpoint,
				blockchain,
				intermediatorSigner,
				intermediatorNonce,
				huobiStorage,
				sr,
			)
			if err != nil {
				return nil, fmt.Errorf("can not create exchange Huobi: (%s)", err.Error())
			}
			exchanges[hb.ID()] = hb
		}
	}
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
