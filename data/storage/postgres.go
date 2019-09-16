package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"

	"github.com/KyberNetwork/reserve-data/common"
	commonv3 "github.com/KyberNetwork/reserve-data/reservesetting/common"
	"github.com/KyberNetwork/reserve-data/world"
)

const (
	schema = `
CREATE TABLE IF NOT EXISTS "fetch_data" 
(
	id SERIAL PRIMARY KEY,
	created TIMESTAMP NOT NULL,
	data JSONB NOT NULL,
	type text NOT NULL
);
CREATE INDEX IF NOT EXISTS "fetch_data_created_index" ON "fetch_data" (created);

CREATE TABLE IF NOT EXISTS "activity"
(
	id SERIAL PRIMARY KEY,
	timepoint BIGINT NOT NULL,
	eid TEXT NOT NULL,
	created TIMESTAMP NOT NULL,
	is_pending BOOL NOT NULL,
	data JSON NOT NULL
);
CREATE INDEX IF NOT EXISTS "activity_idx" ON "activity" (timepoint, eid);

CREATE TABLE IF NOT EXISTS "feed_configuration"
(
	id SERIAL PRIMARY KEY,
	name TEXT UNIQUE NOT NULL,
	enabled BOOLEAN NOT NULL
);
`
	dataTableName          = "fetch_data"
	activityTable          = "activity"
	feedConfigurationTable = "feed_configuration"
	// data type constant

	priceDataType = "price"
	rateDataType  = "rate"
	authDataType  = "authData"
	goldDataType  = "gold"
	btcDataType   = "btc"
)

// PostgresStorage struct
type PostgresStorage struct {
	db *sqlx.DB
}

// NewPostgresStorage return new db instance
func NewPostgresStorage(db *sqlx.DB) (*PostgresStorage, error) {
	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("failed to intialize database schema err=%s", err.Error())
	}

	s := &PostgresStorage{
		db: db,
	}

	// init all feed as enabled
	query := `INSERT INTO "feed_configuration" (name, enabled) VALUES ($1, $2) ON CONFLICT DO NOTHING;`
	for _, feed := range world.AllFeeds() {
		if _, err := s.db.Exec(query, feed, true); err != nil {
			return s, err
		}
	}
	return s, nil
}

func (ps *PostgresStorage) storeData(data interface{}, dataType string, timepoint uint64) error {
	query := fmt.Sprintf(`INSERT INTO "%s" (created, data, type) VALUES ($1, $2, $3)`, dataTableName)
	timestamp := common.TimepointToTime(timepoint)
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if _, err := ps.db.Exec(query, timestamp, dataJSON, dataType); err != nil {
		return err
	}
	return nil
}

func (ps *PostgresStorage) currentVersion(dataType string, timepoint uint64) (common.Version, error) {
	var (
		v  common.Version
		id int64
	)
	timestamp := common.TimepointToTime(timepoint)
	query := fmt.Sprintf(`SELECT id FROM "%s" WHERE created <= $1 and type = $2 ORDER BY created DESC LIMIT 1`, dataTableName)
	if err := ps.db.Get(&id, query, timestamp, dataType); err != nil {
		if err == sql.ErrNoRows {
			return v, fmt.Errorf("there is no version at timestamp: %d", timepoint)
		}
		return v, err
	}
	v = common.Version(id)
	return v, nil
}

// GetData return data from a table
func (ps *PostgresStorage) GetData(table, dataType string, v common.Version) ([]byte, error) {
	var (
		data []byte
	)
	query := fmt.Sprintf(`SELECT data FROM "%s" WHERE id = $1 AND type = $2`, table)
	if err := ps.db.Get(&data, query, v, dataType); err != nil {
		return []byte{}, err
	}
	return data, nil
}

// StorePrice store price
func (ps *PostgresStorage) StorePrice(priceEntry common.AllPriceEntry, timepoint uint64) error {
	return ps.storeData(priceEntry, priceDataType, timepoint)
}

// CurrentPriceVersion return current price version
func (ps *PostgresStorage) CurrentPriceVersion(timepoint uint64) (common.Version, error) {
	return ps.currentVersion(priceDataType, timepoint)
}

// GetAllPrices return all prices currently save in db
func (ps *PostgresStorage) GetAllPrices(v common.Version) (common.AllPriceEntry, error) {
	var (
		allPrices common.AllPriceEntry
	)
	data, err := ps.GetData(dataTableName, priceDataType, v)
	if err != nil {
		return common.AllPriceEntry{}, err
	}
	if err := json.Unmarshal(data, &allPrices); err != nil {
		return common.AllPriceEntry{}, err
	}
	return allPrices, nil
}

// GetOnePrice return one price
func (ps *PostgresStorage) GetOnePrice(pairID uint64, v common.Version) (common.OnePrice, error) {
	allPrices, err := ps.GetAllPrices(v)
	if err != nil {
		return common.OnePrice{}, err
	}
	onePrice, exist := allPrices.Data[pairID]
	if exist {
		return onePrice, nil
	}
	return common.OnePrice{}, errors.New("pair id does not exist")
}

// StoreAuthSnapshot store authdata
func (ps *PostgresStorage) StoreAuthSnapshot(authData *common.AuthDataSnapshot, timepoint uint64) error {
	return ps.storeData(authData, authDataType, timepoint)
}

// CurrentAuthDataVersion return current auth data version
func (ps *PostgresStorage) CurrentAuthDataVersion(timepoint uint64) (common.Version, error) {
	return ps.currentVersion(authDataType, timepoint)
}

// GetAuthData return auth data
func (ps *PostgresStorage) GetAuthData(v common.Version) (common.AuthDataSnapshot, error) {
	var (
		authData common.AuthDataSnapshot
	)
	data, err := ps.GetData(dataTableName, authDataType, v)
	if err != nil {
		return common.AuthDataSnapshot{}, err
	}
	if err := json.Unmarshal(data, &authData); err != nil {
		return common.AuthDataSnapshot{}, err
	}
	return authData, nil
}

// ExportExpiredAuthData export data to store on s3 storage
func (ps *PostgresStorage) ExportExpiredAuthData(timepoint uint64, filePath string) (uint64, error) {
	var (
		data [][]byte
	)
	// Get expire data
	timepointExpireData := timepoint - authDataExpiredDuration
	timestampExpire := common.TimepointToTime(timepointExpireData)
	query := fmt.Sprintf(`SELECT data FROM "%s" WHERE type = $1 AND created > $2`, dataTableName)
	if err := ps.db.Select(&data, query, authDataType, timestampExpire); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	// create export file
	outFile, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer func() {
		if cErr := outFile.Close(); cErr != nil {
			log.Printf("Close file error: %s", cErr.Error())
		}
	}()

	for _, dataByte := range data {
		_, err := outFile.WriteString(string(dataByte) + "\n")
		if err != nil {
			return 0, err
		}
	}
	return uint64(len(data)), nil
}

// PruneExpiredAuthData remove expire auth data from database
func (ps *PostgresStorage) PruneExpiredAuthData(timepoint uint64) (uint64, error) {
	var (
		count uint64
	)
	// Get expire data
	timepointExpireData := timepoint - authDataExpiredDuration
	timestampExpire := common.TimepointToTime(timepointExpireData)
	query := fmt.Sprintf(`WITH deleted AS 
	(DELETE FROM "%s" WHERE type = $1 AND created > $2 RETURNING *) SELECT count(*) FROM deleted`, dataTableName)
	if err := ps.db.Select(&count, query, authDataType, timestampExpire); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return count, nil
}

// StoreRate store rate
func (ps *PostgresStorage) StoreRate(allRateEntry common.AllRateEntry, timepoint uint64) error {
	return ps.storeData(allRateEntry, rateDataType, timepoint)
}

// CurrentRateVersion return current rate version
func (ps *PostgresStorage) CurrentRateVersion(timepoint uint64) (common.Version, error) {
	return ps.currentVersion(rateDataType, timepoint)
}

// GetRate return rate at a specific version
func (ps *PostgresStorage) GetRate(v common.Version) (common.AllRateEntry, error) {
	var (
		rate common.AllRateEntry
	)
	data, err := ps.GetData(dataTableName, rateDataType, v)
	if err != nil {
		return rate, err
	}
	if err := json.Unmarshal(data, &rate); err != nil {
		return common.AllRateEntry{}, err
	}
	return rate, nil
}

//GetRates return rate from time to time
func (ps *PostgresStorage) GetRates(fromTime, toTime uint64) ([]common.AllRateEntry, error) {
	var (
		rates []common.AllRateEntry
		data  [][]byte
	)
	query := fmt.Sprintf(`SELECT data FROM "%s" WHERE type = $1 AND created >= $2 AND created <= $3`, dataTableName)
	from := common.TimepointToTime(fromTime)
	to := common.TimepointToTime(toTime)
	if err := ps.db.Select(&data, query, rateDataType, from, to); err != nil {
		return []common.AllRateEntry{}, err
	}
	for _, dataByte := range data {
		var rate common.AllRateEntry
		if err := json.Unmarshal(dataByte, &rate); err != nil {
			return []common.AllRateEntry{}, err
		}
		rates = append(rates, rate)
	}
	return rates, nil
}

// GetAllRecords return all activities records from database
func (ps *PostgresStorage) GetAllRecords(fromTime, toTime uint64) ([]common.ActivityRecord, error) {
	var (
		activities []common.ActivityRecord
		data       [][]byte
	)
	query := fmt.Sprintf(`SELECT data FROM "%s" WHERE created >= $1 AND created <= $2`, activityTable)
	from := common.TimepointToTime(fromTime)
	to := common.TimepointToTime(toTime)
	if err := ps.db.Select(&data, query, from, to); err != nil {
		return nil, err
	}
	for _, dataByte := range data {
		var activity common.ActivityRecord
		if err := json.Unmarshal(dataByte, &activity); err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}
	return activities, nil
}

// GetPendingActivities return all pending activities
func (ps *PostgresStorage) GetPendingActivities() ([]common.ActivityRecord, error) {
	var (
		pendingActivities []common.ActivityRecord
		data              [][]byte
	)
	query := fmt.Sprintf(`SELECT data FROM "%s" WHERE is_pending IS TRUE`, activityTable)
	if err := ps.db.Select(&data, query); err != nil {
		return []common.ActivityRecord{}, err
	}
	for _, dataByte := range data {
		var activity common.ActivityRecord
		if err := json.Unmarshal(dataByte, &activity); err != nil {
			return []common.ActivityRecord{}, err
		}
		pendingActivities = append(pendingActivities, activity)
	}
	return pendingActivities, nil
}

// UpdateActivity update activity to finished if it is finished
func (ps *PostgresStorage) UpdateActivity(id common.ActivityID, act common.ActivityRecord) error {
	var (
		data []byte
	)
	// get activity from db
	getQuery := fmt.Sprintf(`SELECT data FROM "%s" WHERE timepoint = $1 AND eid = $2`, activityTable)
	if err := ps.db.Get(&data, getQuery, id.Timepoint, id.EID); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	// check if activity is not pending anymore update it
	updateQuery := fmt.Sprintf(`UPDATE "%s" SET is_pending = $1, data = $2 WHERE timepoint = $2 AND eid = $3`, activityTable)
	dataBytes, err := json.Marshal(act)
	if err != nil {
		return err
	}
	if !act.IsPending() {
		if _, err := ps.db.Exec(updateQuery, false, dataBytes, id.Timepoint, id.EID); err != nil {
			return err
		}
	}
	return nil
}

// GetActivity return activity record by id
func (ps *PostgresStorage) GetActivity(id common.ActivityID) (common.ActivityRecord, error) {
	var (
		activityRecord common.ActivityRecord
		data           []byte
	)
	query := fmt.Sprintf(`SELECT data FROM "%s" WHERE timepoint = $1 AND eid = $2`, activityTable)
	if err := ps.db.Get(&data, query, id.Timepoint, id.EID); err != nil {
		return common.ActivityRecord{}, err
	}
	if err := json.Unmarshal(data, &activityRecord); err != nil {
		return common.ActivityRecord{}, err
	}
	return activityRecord, nil
}

// PendingSetRate return pending set rate activity
func (ps *PostgresStorage) PendingSetRate(minedNonce uint64) (*common.ActivityRecord, uint64, error) {
	pendings, err := ps.GetPendingActivities()
	if err != nil {
		return nil, 0, err
	}
	return getFirstAndCountPendingSetrate(pendings, minedNonce)
}

// HasPendingDeposit return true if there is any pending deposit for a token
func (ps *PostgresStorage) HasPendingDeposit(token commonv3.Asset, exchange common.Exchange) (bool, error) {
	var (
		data [][]byte
	)
	query := fmt.Sprintf(`SELECT data FROM "%s" WHERE is_pending IS TRUE AND data->>'Action' = $1`, activityTable)
	if err := ps.db.Select(&data, query, common.ActionDeposit); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	if len(data) > 0 {
		return true, nil
	}
	return false, nil
}

// Record save activity
func (ps *PostgresStorage) Record(action string, id common.ActivityID, destination string,
	params map[string]interface{}, result map[string]interface{},
	estatus string, mstatus string, timepoint uint64) error {
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
	query := fmt.Sprintf(`INSERT INTO "%s" (created, data, is_pending, timepoint, eid) VALUES($1, $2, $3, $4, $5)`, activityTable)
	timestamp := common.TimepointToTime(timepoint)
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}
	if _, err := ps.db.Exec(query, timestamp, data, true, id.Timepoint, id.EID); err != nil {
		return err
	}
	return nil
}

// StoreGoldInfo store gold info into database
func (ps *PostgresStorage) StoreGoldInfo(goldData common.GoldData) error {
	timepoint := goldData.Timestamp
	return ps.storeData(goldData, goldDataType, timepoint)
}

// StoreBTCInfo store btc info into database
func (ps *PostgresStorage) StoreBTCInfo(btcData common.BTCData) error {
	timepoint := btcData.Timestamp
	return ps.storeData(btcData, btcDataType, timepoint)
}

// GetGoldInfo return gold info
func (ps *PostgresStorage) GetGoldInfo(v common.Version) (common.GoldData, error) {
	var (
		goldData common.GoldData
	)
	data, err := ps.GetData(dataTableName, goldDataType, v)
	if err != nil {
		return common.GoldData{}, err
	}
	if err := json.Unmarshal(data, &goldData); err != nil {
		return common.GoldData{}, err
	}
	return goldData, nil
}

// GetBTCInfo return BTC info
func (ps *PostgresStorage) GetBTCInfo(v common.Version) (common.BTCData, error) {
	var (
		btcData common.BTCData
	)
	data, err := ps.GetData(dataTableName, btcDataType, v)
	if err != nil {
		return common.BTCData{}, err
	}
	if err := json.Unmarshal(data, &btcData); err != nil {
		return common.BTCData{}, err
	}
	return btcData, nil
}

// CurrentGoldInfoVersion return btc info version
func (ps *PostgresStorage) CurrentGoldInfoVersion(timepoint uint64) (common.Version, error) {
	return ps.currentVersion(goldDataType, timepoint)
}

// CurrentBTCInfoVersion return current btc info version
func (ps *PostgresStorage) CurrentBTCInfoVersion(timepoint uint64) (common.Version, error) {
	return ps.currentVersion(btcDataType, timepoint)
}

// UpdateFeedConfiguration return false if there is an error
func (ps *PostgresStorage) UpdateFeedConfiguration(name string, enabled bool) error {
	query := fmt.Sprintf(`INSERT INTO %s (name, enabled) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET enabled = EXCLUDED.enabled`, feedConfigurationTable)
	if _, err := ps.db.Exec(query, name, enabled); err != nil {
		return err
	}
	return nil
}

// GetFeedConfiguration return feed configuration
func (ps *PostgresStorage) GetFeedConfiguration() ([]common.FeedConfiguration, error) {
	var (
		result []common.FeedConfiguration
	)
	query := fmt.Sprintf(`SELECT name, enabled FROM "%s"`, feedConfigurationTable)
	if err := ps.db.Select(&result, query); err != nil {
		return nil, err
	}
	return result, nil
}
