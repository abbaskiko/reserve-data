package fetcher

import (
	ethereum "github.com/ethereum/go-ethereum/common"

	"github.com/KyberNetwork/reserve-data/common"
)

// Exchange is the common interface of centralized exchanges.
type Exchange interface {
	ID() common.ExchangeID
	FetchPriceData(timepoint uint64) (map[uint64]common.ExchangePrice, error)
	FetchEBalanceData(timepoint uint64) (common.EBalanceEntry, error)
	OrderStatus(id string, base, quote string) (string, error)
	DepositStatus(id common.ActivityID, txHash string, assetID uint64, amount float64, timepoint uint64) (string, error)
	WithdrawStatus(id string, assetID uint64, amount float64, timepoint uint64) (string, string, error)
	TokenAddresses() (map[string]ethereum.Address, error)
}
