package exchange

import (
	"errors"
	"fmt"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-data/common"
	commonv3 "github.com/KyberNetwork/reserve-data/reservesetting/common"
)

type StableEx struct {
	l *zap.SugaredLogger
}

func (se *StableEx) TokenAddresses() (map[string]ethereum.Address, error) {
	// returning admin multisig. In case anyone sent dgx to this address,
	// we can still get it.
	return map[string]ethereum.Address{
		"DGX":  ethereum.HexToAddress("0xFDF28Bf25779ED4cA74e958d54653260af604C20"),
		"WBTC": ethereum.HexToAddress("0xFDF28Bf25779ED4cA74e958d54653260af604C20"),
	}, nil
}

func (se *StableEx) MarshalText() (text []byte, err error) {
	return []byte(se.ID().String()), nil
}

func (se *StableEx) Address(asset commonv3.Asset) (ethereum.Address, bool) {
	addrs, err := se.TokenAddresses()
	if err != nil {
		return ethereum.Address{}, false
	}
	addr, supported := addrs[asset.Symbol]
	return addr, supported
}

func (se *StableEx) GetLiveExchangeInfos(pairs []commonv3.TradingPairSymbols) (common.ExchangeInfo, error) {
	se.l.Warnw("stabel_exchange shouldn't come with live exchange info. Return an all 0 result...")
	result := make(common.ExchangeInfo)
	for _, pair := range pairs {
		result[pair.ID] = common.ExchangePrecisionLimit{
			Precision:   common.TokenPairPrecision{},
			AmountLimit: common.TokenPairAmountLimit{},
			PriceLimit:  common.TokenPairPriceLimit{},
			MinNotional: 0,
		}
	}
	return result, nil
}

func (se *StableEx) QueryOrder(symbol string, id uint64) (done float64, remaining float64, finished bool, err error) {
	// TODO: see if trade order (a tx to dgx contract) is successful or not
	// - successful: done = order amount, remaining = 0, finished = true, err = nil
	// - failed: done = 0, remaining = order amount, finished = false, err = some error
	// - pending: done = 0, remaining = order amount, finished = false, err = nil
	return 0, 0, false, errors.New("not supported")
}

func (se *StableEx) Trade(tradeType string, pair commonv3.TradingPairSymbols, rate float64, amount float64) (id string, done float64, remaining float64, finished bool, err error) {
	// TODO: communicate with dgx connector to do the trade
	return "not supported", 0, 0, false, errors.New("not supported")
}

func (se *StableEx) Withdraw(asset commonv3.Asset, amount *big.Int, address ethereum.Address) (string, error) {
	// TODO: communicate with dgx connector to withdraw
	return "not supported", errors.New("not supported")
}

func (se *StableEx) CancelOrder(id, base, quote string) error {
	return errors.New("dgx doesn't support trade cancelling")
}

func (se *StableEx) FetchPriceData(timepoint uint64) (map[uint64]common.ExchangePrice, error) {
	result := map[uint64]common.ExchangePrice{}
	// TODO: Get price data from dgx connector and construct valid orderbooks
	return result, nil
}

func (se *StableEx) FetchEBalanceData(timepoint uint64) (common.EBalanceEntry, error) {
	result := common.EBalanceEntry{}
	result.Timestamp = common.Timestamp(fmt.Sprintf("%d", timepoint))
	result.Valid = true
	result.Status = true
	// TODO: Get balance data from dgx connector
	result.ReturnTime = common.GetTimestamp()
	result.AvailableBalance = map[string]float64{"DGX": 0, "ETH": 0}
	result.LockedBalance = map[string]float64{"DGX": 0, "ETH": 0}
	result.DepositBalance = map[string]float64{"DGX": 0, "ETH": 0}
	return result, nil
}

func (se *StableEx) GetTradeHistory(fromTime, toTime uint64) (common.ExchangeTradeHistory, error) {
	return common.ExchangeTradeHistory{}, nil
}

func (se *StableEx) FetchTradeHistory() {
	// TODO: get trade history
}

func (se *StableEx) DepositStatus(id common.ActivityID, txHash string, assetID uint64, amount float64, timepoint uint64) (string, error) {
	// TODO: checking txHash status
	return "", errors.New("not supported")
}

func (se *StableEx) WithdrawStatus(id string, assetID uint64, amount float64, timepoint uint64) (string, string, error) {
	// TODO: checking id (id is the txhash) status
	return "", "", errors.New("not supported")
}

func (se *StableEx) OrderStatus(id string, base, quote string) (string, error) {
	// TODO: checking id (id is the txhash) status
	return "", errors.New("not supported")
}

// Name return exchangeID
func (se *StableEx) ID() common.ExchangeID {
	return common.StableExchange
}

func NewStableEx(l *zap.SugaredLogger) (*StableEx, error) {
	return &StableEx{l: l}, nil
}
