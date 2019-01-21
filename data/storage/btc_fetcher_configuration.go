package storage

import (
	"encoding/json"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/boltdb/bolt"
)

const (
	btcFetcherConfig = "btcFetcherConfig"
)

//UpdateBTCFetcherConfiguration save btc fetcher config to db
func (bs *BoltStorage) UpdateBTCFetcherConfiguration(config common.BTCFetcherConfigurationRequest) (err error) {
	err = bs.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(btcFetcherConfigurationBucket))
		configJSON, uErr := json.Marshal(config.BTC)
		if uErr != nil {
			return uErr
		}
		return b.Put([]byte(btcFetcherConfig), configJSON)
	})
	return err
}

//GetBTCFetcherConfiguration returns config for btc fetcher
func (bs *BoltStorage) GetBTCFetcherConfiguration() (config common.BTCFetcherConfigurationRequest, err error) {
	var (
		configValue bool
	)
	err = bs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(btcFetcherConfigurationBucket))
		v := b.Get([]byte(btcFetcherConfig))
		if v != nil {
			if uErr := json.Unmarshal(v, &configValue); uErr != nil {
				return uErr
			}
			config.BTC = configValue
		}
		return nil
	})
	return config, err
}
