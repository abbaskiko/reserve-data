package core

import (
	"github.com/KyberNetwork/reserve-data/common"
)

// ActivityStorage is the interface contains all database operations of core.
type ActivityStorage interface {
	Record(
		action string,
		id common.ActivityID,
		destination string,
		params map[string]interface{},
		result map[string]interface{},
		estatus string,
		mstatus string,
		timepoint uint64) error
	HasPendingDeposit(
		token common.Token, exchange common.Exchange) (bool, error)

	GetActivity(id common.ActivityID) (common.ActivityRecord, error)

	GetActivityByOrderID(id string) (common.ActivityRecord, error)

	// PendingActivityForAction return the last pending activity and number of pending transactions.
	PendingActivityForAction(minedNonce uint64, activityType string) (*common.ActivityRecord, uint64, error)

	// This activty have been canceled after done
	UpdateCompletedActivity(id common.ActivityID, activity common.ActivityRecord) error
}
