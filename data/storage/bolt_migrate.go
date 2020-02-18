package storage

import (
	"time"

	"github.com/boltdb/bolt"
	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-data/common"
)

// Migrate migrate data to new db
func (bs *BoltStorage) Migrate(newbs *BoltStorage) (err error) {
	defer func() {
		err = newbs.db.Close()
		if err != nil {
			zap.S().Errorw("failed to close dest db", "err", err)
		}
		err = bs.db.Close()
		if err != nil {
			zap.S().Errorw("failed to close source db", "err", err)
		}
	}()
	timeNow := common.TimeToTimepoint(time.Now())
	// btc info
	bs.l.Info("process BTCInfo")
	latestBTCVersion, err := bs.CurrentBTCInfoVersion(timeNow)
	if err == nil {
		latestBTCData, err := bs.GetBTCInfo(latestBTCVersion)
		if err != nil {
			return err
		}
		if err = newbs.StoreBTCInfo(latestBTCData); err != nil {
			return err
		}
	}
	bs.l.Info("process USDInfo")

	// usd info
	latestUSDVersion, err := bs.CurrentUSDInfoVersion(timeNow)
	if err == nil {
		latestUSDData, err := bs.GetUSDInfo(latestUSDVersion)
		if err != nil {
			return err
		}
		if err = newbs.StoreUSDInfo(latestUSDData); err != nil {
			return err
		}
	}

	bs.l.Info("process GoldInfo")

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
	bs.l.Info("process PriceInfo")
	// price info
	latestPriceVersion, err := bs.CurrentPriceVersion(timeNow)
	if err == nil {
		latestPricesData, err := bs.GetAllPrices(latestPriceVersion)
		if err != nil {
			return err
		}
		if err = newbs.StorePrice(latestPricesData, uint64(latestPriceVersion)); err != nil {
			return err
		}
	}

	bs.l.Info("process AuthData")
	// auth data info
	latestAuthDataVersion, err := bs.CurrentAuthDataVersion(timeNow)
	if err == nil {
		latestAuthData, err := bs.GetAuthData(latestAuthDataVersion)
		if err != nil {
			return err
		}
		if err = newbs.StoreAuthSnapshot(&latestAuthData, uint64(latestAuthDataVersion)); err != nil {
			return err
		}
	}

	bs.l.Info("process RateInfo")
	// rate info
	latestRateVersion, err := bs.CurrentRateVersion(timeNow)
	if err == nil {
		latestRateData, err := bs.GetRate(latestRateVersion)
		if err != nil {
			return err
		}
		if err = newbs.StoreRate(latestRateData, uint64(latestRateVersion)); err != nil {
			return err
		}
	}
	reverseForEach := func(b *bolt.Bucket, maxIteration int, fn func(k, v []byte) error) error {
		c := b.Cursor()
		iteration := 0
		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			if err := fn(k, v); err != nil {
				return err
			}
			iteration++
			if iteration == maxIteration {
				break
			}
		}
		return nil
	}

	// other buckets
	errU := newbs.db.Update(func(newTX *bolt.Tx) error {
		errV := bs.db.View(func(tx *bolt.Tx) error {

			bs.l.Info("process metric info")
			mtricBucket := tx.Bucket([]byte(metricBucket))
			mtricBucketNew := newTX.Bucket([]byte(metricBucket))
			maxMetricCopyMax := 10000
			err = reverseForEach(mtricBucket, maxMetricCopyMax, func(k, v []byte) error {
				return mtricBucketNew.Put(k, v)
			})

			for _, bucket := range buckets {
				switch bucket {
				case btcBucket, usdBucket, goldBucket, rateBucket, priceBucket, authDataBucket, metricBucket:
					continue
				default:
					bs.l.Infow("process", "bucket", bucket)
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
