package metric

import (
	"github.com/KyberNetwork/reserve-data/common"
	commonv3 "github.com/KyberNetwork/reserve-data/v3/common"
)

// Storage is the interface that wraps all metrics database operations.
type Storage interface {
	StoreMetric(data *common.MetricEntry, timepoint uint64) error
	StoreRebalanceControl(status bool) error
	StoreSetrateControl(status bool) error

	GetMetric(tokens []commonv3.Asset, fromTime, toTime uint64) (map[uint64]common.MetricList, error)
	GetRebalanceControl() (common.RebalanceControl, error)
	GetSetrateControl() (common.SetrateControl, error)

	SetStableTokenParams(value []byte) error
	ConfirmStableTokenParams(value []byte) error
	RemovePendingStableTokenParams() error
	GetPendingStableTokenParams() (map[string]interface{}, error)
	GetStableTokenParams() (map[string]interface{}, error)
}
