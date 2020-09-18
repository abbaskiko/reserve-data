package data

import (
	"github.com/KyberNetwork/reserve-data/common"
)

// GlobalStorage is the interfaces that wraps database operations of real world
// pricing information of ReserveData.
type GlobalStorage interface {
	GetGoldInfo(version common.Version) (common.GoldData, error)
	CurrentGoldInfoVersion(timepoint uint64) (common.Version, error)

	GetBTCInfo(version common.Version) (common.BTCData, error)
	CurrentBTCInfoVersion(timepoint uint64) (common.Version, error)

	GetUSDInfo(version common.Version) (common.USDData, error)
	CurrentUSDInfoVersion(timepoint uint64) (common.Version, error)

	SetGasThreshold(v common.GasThreshold) error
	GetGasThreshold() (common.GasThreshold, error)

	SetPreferGasSource(v common.PreferGasSource) error
	GetPreferGasSource() (common.PreferGasSource, error)
}
