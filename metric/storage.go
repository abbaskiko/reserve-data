package metric

import (
	"github.com/KyberNetwork/reserve-data/common"
)

// Storage is the interface that wraps all metrics database operations.
type Storage interface {
	StoreMetric(data *common.MetricEntry, timepoint uint64) error
	StoreRebalanceControl(status bool) error
	StoreSetrateControl(status bool) error

	GetMetric(tokens []common.Token, fromTime, toTime uint64) (map[string]common.MetricList, error)
	GetRebalanceControl() (common.RebalanceControl, error)
	GetSetrateControl() (common.SetrateControl, error)

	SetStableTokenParams(value []byte) error
	ConfirmStableTokenParams(value []byte) error
	RemovePendingStableTokenParams() error
	GetPendingStableTokenParams() (map[string]interface{}, error)
	GetStableTokenParams() (map[string]interface{}, error)
}
