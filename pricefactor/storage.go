package pricefactor

import (
	"github.com/KyberNetwork/reserve-data/common"
)

type AssetID uint64

// Storage is the interface that wraps all metrics database operations.
type Storage interface {
	StoreRebalanceControl(status bool) error
	StoreSetrateControl(status bool) error

	GetRebalanceControl() (common.RebalanceControl, error)
	GetSetrateControl() (common.SetrateControl, error)

	SetStableTokenParams(value []byte) error
	ConfirmStableTokenParams(value []byte) error
	RemovePendingStableTokenParams() error
	GetPendingStableTokenParams() (map[string]interface{}, error)
	GetStableTokenParams() (map[string]interface{}, error)
}
