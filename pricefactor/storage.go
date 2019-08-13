package pricefactor

import (
	"github.com/KyberNetwork/reserve-data/common"
	commonv3 "github.com/KyberNetwork/reserve-data/v3/common"
)

type AssetID uint64

// Storage is the interface that wraps all metrics database operations.
type Storage interface {
	StorePriceFactor(data *common.AllPriceFactor, timepoint uint64) error
	StoreRebalanceControl(status bool) error
	StoreSetrateControl(status bool) error

	GetPriceFactor(tokens []commonv3.Asset, fromTime, toTime uint64) (map[AssetID]common.PriceFactorList, error)
	GetRebalanceControl() (common.RebalanceControl, error)
	GetSetrateControl() (common.SetrateControl, error)

	SetStableTokenParams(value []byte) error
	ConfirmStableTokenParams(value []byte) error
	RemovePendingStableTokenParams() error
	GetPendingStableTokenParams() (map[string]interface{}, error)
	GetStableTokenParams() (map[string]interface{}, error)
}
