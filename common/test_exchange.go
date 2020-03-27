package common

import (
	"math/big"

	ethereum "github.com/ethereum/go-ethereum/common"
)

type TestExchange struct {
}

func (te TestExchange) ID() ExchangeID {
	return "binance"
}
func (te TestExchange) Address(token Token) (address ethereum.Address, supported bool) {
	return ethereum.Address{}, true
}
func (te TestExchange) Withdraw(token Token, amount *big.Int, address ethereum.Address, timepoint uint64) (string, error) {
	return "withdrawid", nil
}
func (te TestExchange) Trade(tradeType string, base Token, quote Token, rate float64, amount float64, timepoint uint64) (id string, done float64, remaining float64, finished bool, err error) {
	return "tradeid", 10, 5, false, nil
}
func (te TestExchange) CancelOrder(id, symbol string) error {
	return nil
}
func (te TestExchange) MarshalText() (text []byte, err error) {
	return []byte("bittrex"), nil
}
func (te TestExchange) GetExchangeInfo(pair TokenPairID) (ExchangePrecisionLimit, error) {
	return ExchangePrecisionLimit{}, nil
}
func (te TestExchange) GetFee() (ExchangeFees, error) {
	return ExchangeFees{}, nil
}
func (te TestExchange) GetMinDeposit() (ExchangesMinDeposit, error) {
	return ExchangesMinDeposit{}, nil
}
func (te TestExchange) GetInfo() (ExchangeInfo, error) {
	return ExchangeInfo{}, nil
}
func (te TestExchange) UpdateDepositAddress(token Token, address string) error {
	return nil
}
func (te TestExchange) GetTradeHistory(fromTime, toTime uint64) (ExchangeTradeHistory, error) {
	return ExchangeTradeHistory{}, nil
}

func (te TestExchange) GetLiveExchangeInfos(tokenPairIDs []TokenPairID) (ExchangeInfo, error) {
	return ExchangeInfo{}, nil
}

func (te TestExchange) OpenOrders() ([]Order, error) {
	return nil, nil
}
