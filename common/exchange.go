package common

import (
	"fmt"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum/common"

	"github.com/KyberNetwork/reserve-data/v3/common"
)

// ExchangeName is the name of exchanges of which core will use to rebalance.
//go:generate stringer -type=ExchangeName -linecomment
type ExchangeName int

const (
	//Binance is the enumerated key for binance
	Binance ExchangeName = iota //binance
	//Huobi is the enumerated key for huobi
	Huobi //huobi
	//StableExchange is the enumerated key for stable_exchange
	StableExchange //stable_exchange
)

// ValidExchangeNames returns all valid exchange names.
var ValidExchangeNames = map[string]ExchangeName{
	Binance.String():        Binance,
	Huobi.String():          Huobi,
	StableExchange.String(): StableExchange,
}

// Exchange represents a centralized exchange like Binance, Huobi...
type Exchange interface {
	// TODO ExchangeName should be called ExchangeID, and ExchangeID should be removed
	ID() ExchangeID
	Name() ExchangeName
	// Address return the deposit address of an asset and return true
	// if token is supported in the exchange, otherwise return false.
	// This function will prioritize live address from exchange above the current stored address.
	Address(asset common.Asset) (address ethereum.Address, supported bool)
	Withdraw(asset common.Asset, amount *big.Int, address ethereum.Address, timepoint uint64) (string, error)
	Trade(tradeType string, pair common.TradingPairSymbols, rate, amount float64, timepoint uint64) (id string, done, remaining float64, finished bool, err error)
	CancelOrder(id, base, quote string) error
	MarshalText() (text []byte, err error)

	// GetLiveExchangeInfo querry the Exchange Endpoint for exchange precision and limit of a list of tokenPairIDs
	// It return error if occurs.
	GetLiveExchangeInfos([]common.TradingPairSymbols) (ExchangeInfo, error)
	GetTradeHistory(fromTime, toTime uint64) (ExchangeTradeHistory, error)
}

var SupportedExchanges = map[ExchangeID]Exchange{}

func GetExchange(id string) (Exchange, error) {
	ex := SupportedExchanges[ExchangeID(id)]
	if ex == nil {
		return ex, fmt.Errorf("exchange %s is not supported", id)
	}
	return ex, nil
}

// ExchangeSetting contain the composition of settings necessary for an exchange
// It is use mainly to group all the setting for DB operations
type ExchangeSetting struct {
	DepositAddress ExchangeAddresses   `json:"deposit_address"`
	MinDeposit     ExchangesMinDeposit `json:"min_deposit"`
	Fee            ExchangeFees        `json:"fee"`
	Info           ExchangeInfo        `json:"info"`
}

// NewExchangeSetting returns a pointer to A newly created ExchangeSetting instance
func NewExchangeSetting(depoAddr ExchangeAddresses, minDep ExchangesMinDeposit, fee ExchangeFees, info ExchangeInfo) *ExchangeSetting {
	return &ExchangeSetting{
		DepositAddress: depoAddr,
		MinDeposit:     minDep,
		Fee:            fee,
		Info:           info,
	}
}
