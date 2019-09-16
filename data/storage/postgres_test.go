package storage

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/common/testutil"
)

func TestRate(t *testing.T) {
	db, teardown := testutil.MustNewDevelopmentDB()
	defer func() {
		require.NoError(t, teardown())
	}()

	ps, err := NewPostgresStorage(db)
	require.NoError(t, err)

	// test store data
	baseBuy, _ := big.NewInt(0).SetString("940916409070162411520", 10)
	baseSell, _ := big.NewInt(0).SetString("1051489536265074", 10)
	rateData := common.AllRateEntry{
		Data: map[string]common.RateEntry{
			"KNC": {
				Block:       8539892,
				BaseBuy:     baseBuy,
				BaseSell:    baseSell,
				CompactBuy:  -11,
				CompactSell: 7,
			},
		},
		Timestamp:   "1568358532784",
		ReturnTime:  "1568358532956",
		BlockNumber: 8539899,
	}
	timepoint := uint64(1568358532784)

	// test store rate
	err = ps.StoreData(rateData, dataTableName, rateDataType, timepoint)
	require.NoError(t, err)

	// test get current version
	timepointTest := uint64(1568358532785)
	currentRateVersion, err := ps.CurrentRateVersion(timepointTest)
	require.NoError(t, err)
	assert.Equal(t, common.Version(1), currentRateVersion)

	// test there is no version
	timepointTest = uint64(1568358532783)
	_, err = ps.CurrentRateVersion(timepointTest)
	assert.NotNil(t, err)

	// Test get rate
	rate, err := ps.GetRate(currentRateVersion)
	require.NoError(t, err)
	assert.Equal(t, rateData, rate)
}

func TestPrice(t *testing.T) {
	db, teardown := testutil.MustNewDevelopmentDB()
	defer func() {
		require.NoError(t, teardown())
	}()

	ps, err := NewPostgresStorage(db)
	require.NoError(t, err)

	// test store data
	priceData := common.AllPriceEntry{
		Data: map[uint64]common.OnePrice{
			2: {
				common.ExchangeID(1): common.ExchangePrice{
					Asks: []common.PriceEntry{
						{
							Rate:     0.001062,
							Quantity: 6,
						},
						{
							Rate:     0.0010677,
							Quantity: 376,
						},
					},
					Bids: []common.PriceEntry{
						{
							Rate:     0.0010603,
							Quantity: 46,
						},
						{
							Rate:     0.0010593,
							Quantity: 46,
						},
					},
					Error:      "",
					Valid:      true,
					Timestamp:  "1568358536753",
					ReturnTime: "1568358536834",
				},
			},
		},
		Block: 8539900,
	}

	timepoint := uint64(1568358536753)

	// test store rate
	err = ps.StoreData(priceData, dataTableName, priceDataType, timepoint)
	require.NoError(t, err)

	// test get current version
	timepointTest := uint64(1568358536753)
	currentPriceVersion, err := ps.CurrentPriceVersion(timepointTest)
	require.NoError(t, err)
	assert.Equal(t, common.Version(1), currentPriceVersion)

	// test there is no version
	timepointTest = uint64(1568358532783)
	_, err = ps.CurrentPriceVersion(timepointTest)
	assert.NotNil(t, err)

	// Test get rate
	prices, err := ps.GetAllPrices(currentPriceVersion)
	require.NoError(t, err)
	assert.Equal(t, priceData, prices)
}
