package storage

import (
	"time"

	"github.com/boltdb/bolt"

	"github.com/KyberNetwork/reserve-data/common"
)

// Migrate migrate data to new db
func (bs *BoltStorage) Migrate(newbs *BoltStorage) error {
	timeNow := common.TimeToTimepoint(time.Now())
	// btc info
	latestBTCVersion, err := bs.CurrentBTCInfoVersion(timeNow)
	if err == nil {
		latestBTCData, err := bs.GetBTCInfo(latestBTCVersion)
		if err != nil {
			return err
		}
		if err := newbs.StoreBTCInfo(latestBTCData); err != nil {
			return err
		}
	}

	// usd info
	latestUSDVersion, err := bs.CurrentUSDInfoVersion(timeNow)
	if err == nil {
		latestUSDData, err := bs.GetUSDInfo(latestUSDVersion)
		if err != nil {
			return err
		}
		if err := newbs.StoreUSDInfo(latestUSDData); err != nil {
			return err
		}
	}

	// gold info
	latestGoldVersion, err := bs.CurrentGoldInfoVersion(timeNow)
	if err == nil {
		latestGoldData, err := bs.GetGoldInfo(latestGoldVersion)
		if err != nil {
			return err
		}
		if err := newbs.StoreGoldInfo(latestGoldData); err != nil {
			return err
		}
	}

	// price info
	latestPriceVersion, err := bs.CurrentPriceVersion(timeNow)
	if err == nil {
		latestPricesData, err := bs.GetAllPrices(latestPriceVersion)
		if err != nil {
			return err
		}
		if err := newbs.StorePrice(latestPricesData, uint64(latestPriceVersion)); err != nil {
			return err
		}
	}

	// auth data info
	latestAuthDataVersion, err := bs.CurrentAuthDataVersion(timeNow)
	if err == nil {
		latestAuthData, err := bs.GetAuthData(latestAuthDataVersion)
		if err != nil {
			return err
		}
		if err := newbs.StoreAuthSnapshot(&latestAuthData, uint64(latestAuthDataVersion)); err != nil {
			return err
		}
	}

	// rate info
	latestRateVersion, err := bs.CurrentRateVersion(timeNow)
	if err == nil {
		latestRateData, err := bs.GetRate(latestRateVersion)
		if err != nil {
			return err
		}
		if err := newbs.StoreRate(latestRateData, uint64(latestRateVersion)); err != nil {
			return err
		}
	}

	// other buckets
	errU := newbs.db.Update(func(newTX *bolt.Tx) error {
		errV := bs.db.View(func(tx *bolt.Tx) error {
			for _, bucket := range buckets {
				switch bucket {
				case btcBucket, usdBucket, goldBucket, rateBucket, priceBucket, authDataBucket:
					continue
				default:
					b := tx.Bucket([]byte(bucket))
					newB := newTX.Bucket([]byte(bucket))
					if errFE := b.ForEach(newB.Put); errFE != nil {
						return errFE
					}
				}
			}
			return nil
		})
		return errV
	})

	return errU
}
