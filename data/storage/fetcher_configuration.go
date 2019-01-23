package storage

import (
	"encoding/json"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/boltdb/bolt"
)

const (
	fetcherConfigurationKey = "fetcherConfiguratioKey"
)

//UpdateFetcherConfiguration save btc fetcher config to db
func (bs *BoltStorage) UpdateFetcherConfiguration(config common.FetcherConfiguration) (err error) {
	err = bs.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(fetcherConfigurationBucket))
		configJSON, uErr := json.Marshal(config)
		if uErr != nil {
			return uErr
		}
		if uErr := b.Put([]byte(fetcherConfigurationKey), configJSON); uErr != nil {
			return uErr
		}
		return nil
	})
	return err
}

//GetAllFetcherConfiguration returns config for fetcher
func (bs *BoltStorage) GetAllFetcherConfiguration() (config common.FetcherConfiguration, err error) {
	err = bs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(fetcherConfigurationBucket))
		v := b.Get([]byte(fetcherConfigurationKey))
		if v != nil {
			if uErr := json.Unmarshal(v, &config); uErr != nil {
				return uErr
			}
		}
		return nil
	})
	return config, err
}
