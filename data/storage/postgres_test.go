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
		Data: map[uint64]common.RateEntry{
			2: {
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
	err = ps.StoreRate(rateData, timepoint)
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

	// test store price
	err = ps.StorePrice(priceData, timepoint)
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

func TestActivity(t *testing.T) {
	db, teardown := testutil.MustNewDevelopmentDB()
	defer func() {
		require.NoError(t, teardown())
	}()

	ps, err := NewPostgresStorage(db)
	require.NoError(t, err)

	activityTest := common.ActivityRecord{
		Action: "deposit",
		ID: common.ActivityID{
			Timepoint: 1568622132671609009,
			EID:       "0x7437e2ac582a7cdef75a6c8355d03167a8ab7670a178197d81f14cea76684d74|BQX|39811.443679",
		},
		Destination: "binance",
		Params: common.ActivityParams{
			Amount:    39811.443679,
			Exchange:  common.Binance,
			Timepoint: uint64(1568622125860),
			Asset:     2, // KNC id
		},
		Result: common.ActivityResult{
			BlockNumber: 8559409,
			Error:       "",
			GasPrice:    "50100000000",
			Nonce:       "11039",
			StatusError: "",
			Tx:          "0x7437e2ac582a7cdef75a6c8355d03167a8ab7670a178197d81f14cea76684d74",
		},
		ExchangeStatus: "",
		MiningStatus:   "mined",
		Timestamp:      "1568622125860",
	}
	err = ps.Record(activityTest.Action, activityTest.ID, activityTest.Destination,
		activityTest.Params, activityTest.Result, activityTest.ExchangeStatus, activityTest.MiningStatus, 1568622125860)
	assert.NoError(t, err)

	// test update activity
	testID := common.ActivityID{
		Timepoint: 1568622132671609009,
		EID:       "0x7437e2ac582a7cdef75a6c8355d03167a8ab7670a178197d81f14cea76684d74|BQX|39811.443679",
	}
	
	activityTest.ExchangeStatus = common.ExchangeStatusDone
	err = ps.UpdateActivity(testID, activityTest)
	assert.NoError(t, err)

	// test get activity
	activity, err := ps.GetActivity(testID)
	assert.NoError(t, err)
	assert.Equal(t, activityTest, activity)
}
