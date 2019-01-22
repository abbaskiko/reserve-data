package storage

import (
	"encoding/json"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/boltdb/bolt"
)

//UpdateFetcherConfiguration save btc fetcher config to db
func (bs *BoltStorage) UpdateFetcherConfiguration(config common.FetcherConfiguration) (err error) {
	err = bs.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(btcFetcherConfigurationBucket))
		for k, v := range config.Config {

			configJSON, uErr := json.Marshal(v)
			if uErr != nil {
				return uErr
			}
			if uErr := b.Put([]byte(k), configJSON); uErr != nil {
				return uErr
			}
		}
		return nil
	})
	return err
}

//GetAllFetcherConfiguration returns config for fetcher
func (bs *BoltStorage) GetAllFetcherConfiguration() (config common.FetcherConfiguration, err error) {
	var (
		configValue bool
	)
	config.Config = make(map[string]bool)
	err = bs.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(btcFetcherConfigurationBucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if v != nil {
				if uErr := json.Unmarshal(v, &configValue); uErr != nil {
					return uErr
				}
				config.Config[string(k)] = configValue
			}
		}
		return nil
	})
	return config, err
}
