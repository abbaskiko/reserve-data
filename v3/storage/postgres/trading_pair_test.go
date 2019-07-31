package postgres

import (
	"testing"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/common/testutil"
	commonv3 "github.com/KyberNetwork/reserve-data/v3/common"
	"github.com/KyberNetwork/reserve-data/v3/storage"
)

func TestStorage_CreateTradingPair(t *testing.T) {
	db, tearDown := testutil.MustNewDevelopmentDB()
	defer func() {
		assert.NoError(t, tearDown())
	}()

	s, err := NewStorage(db)
	require.NoError(t, err)

	baseAssetID, err := s.CreateAsset(
		"ABC",
		"ABC Advanced Token",
		ethereum.HexToAddress("0xB42b3F0C10385df057f7374d7Aa884571E71791b"),
		12,
		false,
		commonv3.SetRateNotSet,
		false,
		true,
		nil,
		nil,
		[]commonv3.AssetExchange{
			{
				ExchangeID:     uint64(common.Binance),
				Symbol:         "ABC",
				DepositAddress: ethereum.HexToAddress("0x118ee757dD8841F81903E1C1d7d7Aa88e376cC39"),
				MinDeposit:     1,
				WithdrawFee:    2,
				TradingPairs:   nil,
			},
		},
		testAssetTarget,
	)
	require.NoError(t, err)

	baseAssetID2, err := s.CreateAsset(
		"DEF",
		"DEF Advanced Token",
		ethereum.HexToAddress("0x464b0b37db1ee1b5fbe27300acfbf172fd5e4f53"),
		12,
		false,
		commonv3.SetRateNotSet,
		false,
		true,
		nil,
		nil,
		[]commonv3.AssetExchange{},
		testAssetTarget,
	)

	_, err = s.CreateAssetExchange(
		uint64(common.Binance),
		1, // ETH ID
		"ETH",
		ethereum.HexToAddress("0x7523b6251F18EB323EEEF004662DE50e1Cedfd95"),
		0.0001,
		0.0001,
		0.001,
		0.001,
	)
	require.NoError(t, err)

	var tests = []struct {
		msg             string
		exchangeID      uint64
		baseID          uint64
		quoteID         uint64
		pricePrecision  uint64
		amountPrecision uint64
		amountLimitMin  float64
		amountLimitMax  float64
		priceLimitMin   float64
		priceLimitMax   float64
		minNotional     float64
		assertFn        func(*testing.T, uint64, error)
	}{
		{
			msg:             "create trading pair success",
			exchangeID:      uint64(common.Binance),
			baseID:          baseAssetID,
			quoteID:         1,
			pricePrecision:  6,
			amountPrecision: 10,
			amountLimitMin:  1000,
			amountLimitMax:  10000,
			priceLimitMin:   0.1,
			priceLimitMax:   10.10,
			minNotional:     0.001,
			assertFn: func(t *testing.T, id uint64, err error) {
				require.NoError(t, err)
				assert.NotZero(t, id)
			},
		},
		{
			msg:             "create trading pair with invalid constraint UNIQUE (exchange_id, base_id, quote_id)",
			exchangeID:      uint64(common.Binance),
			baseID:          baseAssetID,
			quoteID:         1,
			pricePrecision:  6,
			amountPrecision: 10,
			amountLimitMin:  1000,
			amountLimitMax:  10000,
			priceLimitMin:   0.1,
			priceLimitMax:   10.10,
			minNotional:     0.001,
			assertFn: func(t *testing.T, id uint64, err error) {
				require.NotNil(t, err)
				// require.NotEqual(t, commonv3.ErrBadTradingPairConfiguration, err)
				assert.Zero(t, id)
			},
		},
		{
			msg:             "create trading pair with inexistent exchange id",
			exchangeID:      1234,
			baseID:          baseAssetID,
			quoteID:         1,
			pricePrecision:  6,
			amountPrecision: 10,
			amountLimitMin:  1000,
			amountLimitMax:  10000,
			priceLimitMin:   0.1,
			priceLimitMax:   10.10,
			minNotional:     0.001,
			assertFn: func(t *testing.T, id uint64, err error) {
				require.Zero(t, id)
				require.NotNil(t, err)
				// require.Equal(t, commonv3.ErrBadTradingPairConfiguration, err)
			},
		},
		{
			msg:             "create trading pair with inexistent quote id",
			exchangeID:      uint64(common.Binance),
			baseID:          baseAssetID,
			quoteID:         1234,
			pricePrecision:  6,
			amountPrecision: 10,
			amountLimitMin:  1000,
			amountLimitMax:  10000,
			priceLimitMin:   0.1,
			priceLimitMax:   10.10,
			minNotional:     0.001,
			assertFn: func(t *testing.T, id uint64, err error) {
				require.Zero(t, id)
				require.NotNil(t, err)
				// require.Equal(t, commonv3.ErrBadTradingPairConfiguration, err)
			},
		},
		{
			msg:             "create trading pair with inexistent base id",
			exchangeID:      uint64(common.Binance),
			baseID:          1234,
			quoteID:         1,
			pricePrecision:  6,
			amountPrecision: 10,
			amountLimitMin:  1000,
			amountLimitMax:  10000,
			priceLimitMin:   0.1,
			priceLimitMax:   10.10,
			minNotional:     0.001,
			assertFn: func(t *testing.T, id uint64, err error) {
				require.Zero(t, id)
				require.NotNil(t, err)
				// require.Equal(t, commonv3.ErrBadTradingPairConfiguration, err)
			},
		},
		{
			msg:             "create trading pair exist base id but without asset exchange",
			exchangeID:      uint64(common.Binance),
			baseID:          baseAssetID2,
			quoteID:         1,
			pricePrecision:  6,
			amountPrecision: 10,
			amountLimitMin:  1000,
			amountLimitMax:  10000,
			priceLimitMin:   0.1,
			priceLimitMax:   10.10,
			minNotional:     0.001,
			assertFn: func(t *testing.T, id uint64, err error) {
				require.Zero(t, id)
				require.NotNil(t, err)
			},
		},
	}

	tx, err := s.db.Beginx()

	defer rollbackUnlessCommitted(tx)
	require.NoError(t, err)

	for _, tc := range tests {
		t.Logf("run test case: %s", tc.msg)
		testTPID, err := s.createTradingPair(
			tx,
			tc.exchangeID,
			tc.baseID,
			tc.quoteID,
			tc.pricePrecision,
			tc.amountPrecision,
			tc.amountLimitMin,
			tc.amountLimitMax,
			tc.priceLimitMin,
			tc.priceLimitMax,
			tc.minNotional,
		)
		tc.assertFn(t, testTPID, err)
	}
}

func TestStorage_UpdateTradingPair(t *testing.T) {
	db, tearDown := testutil.MustNewDevelopmentDB()
	defer func() {
		assert.NoError(t, tearDown())
	}()

	s, err := NewStorage(db)
	require.NoError(t, err)

	quoteAssetID, err := s.CreateAsset(
		"ABC",
		"ABC Advanced Token",
		ethereum.HexToAddress("0xB42b3F0C10385df057f7374d7Aa884571E71791b"),
		12,
		false,
		commonv3.SetRateNotSet,
		false,
		true,
		nil,
		nil,
		[]commonv3.AssetExchange{
			{
				ExchangeID:     uint64(common.Binance),
				Symbol:         "ABC",
				DepositAddress: ethereum.HexToAddress("0x118ee757dD8841F81903E1C1d7d7Aa88e376cC39"),
				MinDeposit:     1,
				WithdrawFee:    2,
				TradingPairs:   nil,
			},
		},
		testAssetTarget,
	)
	require.NoError(t, err)

	_, err = s.CreateAsset(
		"DEF",
		"DEF Super Token",
		ethereum.HexToAddress("0xffe97fe10290715ba416a7c1Fd265F28dc574dd9"),
		12,
		false,
		commonv3.ExchangeFeed,
		true,
		false,
		testPWI,
		testRb,
		[]commonv3.AssetExchange{
			{
				ExchangeID:     uint64(common.Binance),
				Symbol:         "BNB",
				DepositAddress: ethereum.HexToAddress("0x118ee757dD8841F81903E1C1d7d7Aa88e376cC39"),
				MinDeposit:     1,
				WithdrawFee:    2,
				TradingPairs: []commonv3.TradingPair{
					{
						Base:            0,
						Quote:           quoteAssetID, // ETH
						PricePrecision:  1,
						AmountPrecision: 2,
						AmountLimitMin:  3.0,
						AmountLimitMax:  4.0,
						PriceLimitMin:   5.0,
						PriceLimitMax:   6.0,
						MinNotional:     1.2,
					},
				},
			},
		},
		testAssetTarget,
	)
	require.NoError(t, err)

	pairs, err := s.GetTradingPairs(uint64(common.Binance))
	require.NoError(t, err)
	require.Len(t, pairs, 1)
	pairID := pairs[0].ID
	var (
		pricePrecision  = uint64(100)
		amountPrecision = uint64(200)
		amountLimitMin  = 1.3
		amountLimitMax  = 1.4
		priceLimitMin   = 1.5
		priceLimitMax   = 1.6
		minNotional     = 1.7
	)
	err = s.UpdateTradingPair(pairID, storage.UpdateTradingPairOpts{
		PricePrecision:  &pricePrecision,
		AmountPrecision: &amountPrecision,
		AmountLimitMin:  &amountLimitMin,
		AmountLimitMax:  &amountLimitMax,
		PriceLimitMin:   &priceLimitMin,
		PriceLimitMax:   &priceLimitMax,
		MinNotional:     &minNotional,
	})
	require.NoError(t, err)
	pairs, err = s.GetTradingPairs(uint64(common.Binance))
	require.NoError(t, err)
	require.Len(t, pairs, 1)
	pair := pairs[0]
	assert.Equal(t, uint64(100), pair.PricePrecision)
	assert.Equal(t, uint64(200), pair.AmountPrecision)
	assert.Equal(t, 1.3, pair.AmountLimitMin)
	assert.Equal(t, 1.4, pair.AmountLimitMax)
	assert.Equal(t, 1.5, pair.PriceLimitMin)
	assert.Equal(t, 1.6, pair.PriceLimitMax)
	assert.Equal(t, 1.7, pair.MinNotional)
}
