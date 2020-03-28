package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

//Record save activity
func (bs *BoltStorage) Record(
	action string,
	id common.ActivityID,
	destination string,
	params map[string]interface{}, result map[string]interface{},
	estatus string,
	mstatus string,
	timepoint uint64) error {

	var err error
	err = bs.db.Update(func(tx *bolt.Tx) error {
		var dataJSON []byte
		b := tx.Bucket([]byte(activityBucket))
		record := common.NewActivityRecord(
			action,
			id,
			destination,
			params,
			result,
			estatus,
			mstatus,
			common.Timestamp(strconv.FormatUint(timepoint, 10)),
		)
		dataJSON, err = json.Marshal(record)
		if err != nil {
			return err
		}

		idByte := id.ToBytes()
		err = b.Put(idByte[:], dataJSON)
		if err != nil {
			return err
		}
		if record.IsPending() {
			pb := tx.Bucket([]byte(pendingActivityBucket))
			// all other pending set rates should be staled now
			// remove all of them
			// AFTER EXPERIMENT, THIS WILL NOT WORK
			// log.Printf("===> Trying to remove staled set rates")
			// if record.Action == "set_rates" {
			// 	stales := []common.ActivityRecord{}
			// 	c := pb.Cursor()
			// 	for k, v := c.First(); k != nil; k, v = c.Next() {
			// 		record := common.ActivityRecord{}
			// 		log.Printf("===> staled act: %+v", record)
			// 		err = json.Unmarshal(v, &record)
			// 		if err != nil {
			// 			return err
			// 		}
			// 		if record.Action == "set_rates" {
			// 			stales = append(stales, record)
			// 		}
			// 	}
			// 	log.Printf("===> removing staled acts: %+v", stales)
			// 	bs.RemoveStalePendingActivities(tx, stales)
			// }
			// after remove all of them, put new set rate activity
			err = pb.Put(idByte[:], dataJSON)
		}
		return err
	})
	return err
}

func formatTimepointToActivityID(timepoint uint64, id []byte) []byte {
	if timepoint == 0 {
		return id
	}
	activityID := common.NewActivityID(timepoint, "")
	byteID := activityID.ToBytes()
	return byteID[:]
}

//GetActivity get activity
func (bs *BoltStorage) GetActivity(id common.ActivityID) (common.ActivityRecord, error) {
	result := common.ActivityRecord{}

	err := bs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(activityBucket))
		idBytes := id.ToBytes()
		v := b.Get(idBytes[:])
		if v == nil {
			return errors.New("can not find that activity")
		}
		return json.Unmarshal(v, &result)
	})
	return result, err
}

// GetActivityByOrderID return activity by eid
func (bs *BoltStorage) GetActivityByOrderID(id string) (common.ActivityRecord, error) {
	result := common.ActivityRecord{}
	var err error
	err = bs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(activityBucket))
		c := b.Cursor()
		fKey, _ := c.First()
		lKey, _ := c.Last()
		toTime := common.TimeToTimepoint(time.Now())
		fromTime := toTime - maxSearchRange // one day from now
		min := formatTimepointToActivityID(fromTime, fKey)
		max := formatTimepointToActivityID(toTime, lKey)
		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			record := common.ActivityRecord{}
			err = json.Unmarshal(v, &record)
			if record.ID.EID == id {
				result = record
				break
			}
		}
		return nil
	})
	return result, err
}

// GetAllRecords return all activity records
func (bs *BoltStorage) GetAllRecords(fromTime, toTime uint64) ([]common.ActivityRecord, error) {
	result := []common.ActivityRecord{}
	var err error
	if (toTime-fromTime)/1000000 > maxGetRatesPeriod {
		return result, fmt.Errorf("time range is too broad, it must be smaller or equal to %d miliseconds", maxGetRatesPeriod)
	}
	err = bs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(activityBucket))
		c := b.Cursor()
		fkey, _ := c.First()
		lkey, _ := c.Last()
		min := formatTimepointToActivityID(fromTime, fkey)
		max := formatTimepointToActivityID(toTime, lkey)

		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			record := common.ActivityRecord{}
			err = json.Unmarshal(v, &record)
			if err != nil {
				return err
			}
			result = append([]common.ActivityRecord{record}, result...)
		}
		return err
	})
	return result, err
}

func getFirstAndCountPendingAction(l *zap.SugaredLogger, pendings []common.ActivityRecord, minedNonce uint64,
	activityType string) (*common.ActivityRecord, uint64, error) {
	var minNonce uint64 = math.MaxUint64
	var minPrice uint64 = math.MaxUint64
	var result *common.ActivityRecord
	var count uint64
	for i, act := range pendings {
		if act.Action == activityType {
			l.Infof("looking for pending (%s): %+v", activityType, act)
			nonce := interfaceConverstionToUint64(l, act.Result["nonce"])
			if nonce < minedNonce {
				l.Infof("NONCE_ISSUE: stalled pending %s transaction, pending: %d, mined: %d",
					activityType, nonce, minedNonce)
				continue
			} else if nonce-minedNonce > 1 {
				l.Infof("NONCE_ISSUE: pending %s transaction for inconsecutive nonce, mined nonce: %d, request nonce: %d",
					activityType, minedNonce, nonce)
			}

			gasPrice := interfaceConverstionToUint64(l, act.Result["gasPrice"])
			if nonce == minNonce {
				if gasPrice < minPrice {
					minNonce = nonce
					result = &pendings[i]
					minPrice = gasPrice
				}
				count++
			} else if nonce < minNonce {
				minNonce = nonce
				result = &pendings[i]
				minPrice = gasPrice
				count = 1
			}
		}
	}

	if result == nil {
		l.Infof("NONCE_ISSUE: found no pending %s transaction with nonce newer than equal to mined nonce: %d",
			activityType, minedNonce)
	} else {
		l.Infof("NONCE_ISSUE: un-mined pending %s, nonce: %d, count: %d, mined nonce: %d",
			activityType, interfaceConverstionToUint64(l, result.Result["nonce"]), count, minedNonce)
	}

	return result, count, nil
}

//RemoveStalePendingActivities remove it
func (bs *BoltStorage) RemoveStalePendingActivities(tx *bolt.Tx, stales []common.ActivityRecord) error {
	pb := tx.Bucket([]byte(pendingActivityBucket))
	for _, stale := range stales {
		idBytes := stale.ID.ToBytes()
		if err := pb.Delete(idBytes[:]); err != nil {
			return err
		}
	}
	return nil
}

//PendingActivityForAction return pending set rate activity
func (bs *BoltStorage) PendingActivityForAction(minedNonce uint64, activityType string) (ar *common.ActivityRecord, count uint64, err error) {
	pendings, err := bs.GetPendingActivities()
	if err != nil {
		return nil, 0, err
	}
	return getFirstAndCountPendingAction(bs.l, pendings, minedNonce, activityType)
}

//GetPendingActivities return pending activities
func (bs *BoltStorage) GetPendingActivities() ([]common.ActivityRecord, error) {
	result := []common.ActivityRecord{}
	var err error
	err = bs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(pendingActivityBucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			record := common.ActivityRecord{}
			err = json.Unmarshal(v, &record)
			if err != nil {
				return err
			}
			result = append(
				[]common.ActivityRecord{record}, result...)
		}
		return err
	})
	return result, err
}

//UpdateActivity update activity info
func (bs *BoltStorage) UpdateActivity(id common.ActivityID, activity common.ActivityRecord) error {
	err := bs.db.Update(func(tx *bolt.Tx) error {
		pb := tx.Bucket([]byte(pendingActivityBucket))
		idBytes := id.ToBytes()
		dataJSON, uErr := json.Marshal(activity)
		if uErr != nil {
			return uErr
		}
		// only update when it exists in pending activity bucket because
		// It might be deleted if it is replaced by another activity
		found := pb.Get(idBytes[:])
		if found != nil {
			uErr = pb.Put(idBytes[:], dataJSON)
			if uErr != nil {
				return uErr
			}
			if !activity.IsPending() {
				uErr = pb.Delete(idBytes[:])
				if uErr != nil {
					return uErr
				}
			}
		}
		b := tx.Bucket([]byte(activityBucket))
		return b.Put(idBytes[:], dataJSON)
	})
	return err
}

//HasPendingDeposit check if a deposit is pending
func (bs *BoltStorage) HasPendingDeposit(token common.Token, exchange common.Exchange) (bool, error) {
	var (
		err    error
		result = false
	)
	err = bs.db.View(func(tx *bolt.Tx) error {
		pb := tx.Bucket([]byte(pendingActivityBucket))
		c := pb.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			record := common.ActivityRecord{}
			if uErr := json.Unmarshal(v, &record); uErr != nil {
				return uErr
			}
			if record.Action == common.ActionDeposit {
				tokenID, ok := record.Params["token"].(string)
				if !ok {
					bs.l.Warnw("record Params token can not be converted to string", "token", record.Params["token"])
					continue
				}
				if tokenID == token.ID && record.Destination == string(exchange.ID()) {
					result = true
				}
			}
		}
		return nil
	})
	return result, err
}

// UpdateCompletedActivity update completed activity to canceled
func (bs *BoltStorage) UpdateCompletedActivity(id common.ActivityID, activity common.ActivityRecord) error {
	err := bs.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(activityBucket))
		idBytes := id.ToBytes()
		found := b.Get(idBytes[:])
		if found != nil {
			dataJSON, uErr := json.Marshal(activity)
			if uErr != nil {
				return uErr
			}
			return b.Put(idBytes[:], dataJSON)
		}
		return nil
	})
	return err
}
