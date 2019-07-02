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

var (
	testPWI = &commonv3.AssetPWI{
		Ask: commonv3.PWIEquation{
			A:                   1,
			B:                   2,
			C:                   3,
			MinMinSpread:        4,
			PriceMultiplyFactor: 5,
		},
		Bid: commonv3.PWIEquation{
			A:                   6,
			B:                   7,
			C:                   8,
			MinMinSpread:        9,
			PriceMultiplyFactor: 10,
		},
	}
	testRb = &commonv3.RebalanceQuadratic{
		A: 1,
		B: 2,
		C: 3,
	}
	testAssetExchanges = []commonv3.AssetExchange{
		{
			ExchangeID:        uint64(common.Binance),
			Symbol:            "BNB",
			DepositAddress:    ethereum.HexToAddress("0x118ee757dD8841F81903E1C1d7d7Aa88e376cC39"),
			MinDeposit:        1,
			WithdrawFee:       2,
			PricePrecision:    3,
			AmountPrecision:   4,
			AmountLimitMin:    5,
			AmountLimitMax:    6,
			PriceLimitMin:     7,
			PriceLimitMax:     8,
			TargetRecommended: 12.0,
			TargetRatio:       12.1,
		},
		{
			ExchangeID:        uint64(common.Huobi),
			Symbol:            "HUO",
			DepositAddress:    ethereum.HexToAddress("0x71241678e935f07ff78182F41881214B77d8cD99"),
			MinDeposit:        2,
			WithdrawFee:       3,
			PricePrecision:    4,
			AmountPrecision:   5,
			AmountLimitMin:    6,
			AmountLimitMax:    7,
			PriceLimitMin:     8,
			PriceLimitMax:     9,
			TargetRecommended: 13.0,
			TargetRatio:       13.1,
		},
	}
	testAssetTarget = &commonv3.AssetTarget{
		Total:              123.1,
		Reserve:            50.2,
		RebalanceThreshold: 50.3,
		TransferThreshold:  50.4,
	}
)

func TestStorage_CreateAsset(t *testing.T) {
	db, tearDown := testutil.MustNewDevelopmentDB()
	defer func() {
		assert.NoError(t, tearDown())
	}()

	s, err := NewStorage(db)
	require.NoError(t, err)

	var tests = []struct {
		msg string

		// input
		symbol    string
		name      string
		address   ethereum.Address
		decimals  uint64
		setRate   commonv3.SetRate
		rebalance bool
		isQuote   bool
		pwi       *commonv3.AssetPWI
		rb        *commonv3.RebalanceQuadratic
		exchanges []commonv3.AssetExchange
		target    *commonv3.AssetTarget

		// output
		assertFn func(*testing.T, uint64, error)
	}{
		{
			msg:       "creating asset without address for non quote token",
			symbol:    "ABC",
			name:      "ABC Advanced Token",
			decimals:  12,
			setRate:   commonv3.BTCFeed,
			rebalance: true,
			isQuote:   false,
			pwi:       testPWI,
			rb:        testRb,
			exchanges: testAssetExchanges,
			target:    testAssetTarget,
			assertFn: func(t *testing.T, id uint64, err error) {
				require.Equal(t, commonv3.ErrAddressMissing, err)
			},
		},
		{
			msg:       "creating asset successfully",
			symbol:    "ABC",
			name:      "ABC Advanced Token",
			address:   ethereum.HexToAddress("0xB8c77482e45F1F44dE1745F52C74426C631bDD52"),
			decimals:  12,
			setRate:   commonv3.BTCFeed,
			rebalance: true,
			isQuote:   false,
			pwi:       testPWI,
			rb:        testRb,
			exchanges: testAssetExchanges,
			target:    testAssetTarget,
			assertFn: func(t *testing.T, id uint64, err error) {
				require.NoError(t, err)
				assert.NotZero(t, id)
			},
		},
		{
			msg:       "creating asset with duplicated symbol",
			symbol:    "ABC",
			name:      "ABC Advanced Token 2",
			address:   ethereum.HexToAddress("0xD2b6Ba1e59373A2750F3D9fE9178706fBd42F1F2"),
			decimals:  12,
			setRate:   commonv3.BTCFeed,
			rebalance: true,
			isQuote:   false,
			pwi:       testPWI,
			rb:        testRb,
			exchanges: testAssetExchanges,
			target:    testAssetTarget,
			assertFn: func(t *testing.T, id uint64, err error) {
				require.Equal(t, commonv3.ErrSymbolExists, err)
			},
		},
		{
			msg:       "creating order with duplicated address",
			symbol:    "ABC-2",
			name:      "ABC Advanced Token-2",
			address:   ethereum.HexToAddress("0xB8c77482e45F1F44dE1745F52C74426C631bDD52"),
			decimals:  12,
			setRate:   commonv3.GoldFeed,
			rebalance: false,
			isQuote:   false,
			pwi:       testPWI,
			rb:        nil,
			exchanges: nil,
			target:    nil,
			assertFn: func(t *testing.T, id uint64, err error) {
				require.Equal(t, commonv3.ErrAddressExists, err)
			},
		},
		{
			msg:       "creating a quote token with null address",
			symbol:    "BTCXYZ",
			name:      "Bitcoin Fork XYZ",
			decimals:  12,
			setRate:   commonv3.SetRateNotSet,
			rebalance: false,
			isQuote:   true,
			rb:        nil,
			exchanges: nil,
			target:    nil,
			assertFn: func(t *testing.T, id uint64, err error) {
				require.NoError(t, err)
			},
		},
		{
			msg:       "creating asset with set rate strategy but no pwi configuration",
			symbol:    "Dodge Coin",
			name:      "Barf",
			address:   ethereum.HexToAddress("0xa57E3c6A7A1A2f5834f41b6B9545d5591dBcE8E0"),
			decimals:  9,
			setRate:   commonv3.ExchangeFeed,
			rebalance: true,
			isQuote:   false,
			rb:        testRb,
			exchanges: testAssetExchanges,
			target:    testAssetTarget,
			assertFn: func(t *testing.T, id uint64, err error) {
				require.Equal(t, commonv3.ErrPWIMissing, err)
			},
		},
		{
			msg:       "creating asset with rebalance but no rebalance quadratic",
			symbol:    "Dodge Coin",
			name:      "Barf",
			address:   ethereum.HexToAddress("0xa57E3c6A7A1A2f5834f41b6B9545d5591dBcE8E0"),
			decimals:  9,
			setRate:   commonv3.ExchangeFeed,
			rebalance: true,
			isQuote:   false,
			pwi:       testPWI,
			rb:        nil,
			exchanges: testAssetExchanges,
			target:    testAssetTarget,
			assertFn: func(t *testing.T, id uint64, err error) {
				require.Equal(t, commonv3.ErrRebalanceQuadraticMissing, err)
			},
		},
		{
			msg:       "creating asset with rebalance but no exchange configuration",
			symbol:    "Dodge Coin",
			name:      "Barf",
			address:   ethereum.HexToAddress("0xa57E3c6A7A1A2f5834f41b6B9545d5591dBcE8E0"),
			decimals:  9,
			setRate:   commonv3.ExchangeFeed,
			rebalance: true,
			isQuote:   false,
			pwi:       testPWI,
			rb:        testRb,
			exchanges: nil,
			target:    testAssetTarget,
			assertFn: func(t *testing.T, id uint64, err error) {
				require.Equal(t, commonv3.ErrAssetExchangeMissing, err)
			},
		},
		{
			msg:       "creating asset with rebalance but no target configuration",
			symbol:    "Dodge Coin",
			name:      "Barf",
			address:   ethereum.HexToAddress("0xa57E3c6A7A1A2f5834f41b6B9545d5591dBcE8E0"),
			decimals:  9,
			setRate:   commonv3.ExchangeFeed,
			rebalance: true,
			isQuote:   false,
			pwi:       testPWI,
			rb:        testRb,
			exchanges: testAssetExchanges,
			target:    nil,
			assertFn: func(t *testing.T, id uint64, err error) {
				require.Equal(t, commonv3.ErrAssetTargetMissing, err)
			},
		},
	}

	for _, tc := range tests {
		t.Logf("running test case for: %s", tc.msg)
		id, err := s.CreateAsset(
			tc.symbol,
			tc.name,
			tc.address,
			tc.decimals,
			tc.setRate,
			tc.rebalance,
			tc.isQuote,
			tc.pwi,
			tc.rb,
			tc.exchanges,
			tc.target,
		)
		tc.assertFn(t, id, err)
	}
}

func TestStorage_GetAssets(t *testing.T) {
	db, tearDown := testutil.MustNewDevelopmentDB()
	defer func() {
		assert.NoError(t, tearDown())
	}()

	s, err := NewStorage(db)
	require.NoError(t, err)

	assets, err := s.GetAssets()
	require.NoError(t, err)
	assert.Len(t, assets, 1, "expected that ETH is initialized")
	ethAsset := assets[0]
	assert.Equal(t, uint64(1), ethAsset.ID)
	assert.Equal(t, "ETH", ethAsset.Symbol)
	assert.Equal(t, "Ethereum", ethAsset.Name)
	assert.Equal(t, ethereum.HexToAddress("0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"), ethAsset.Address)
	assert.Equal(t, uint64(18), ethAsset.Decimals)
	assert.Equal(t, commonv3.SetRateNotSet, ethAsset.SetRate)
	assert.False(t, ethAsset.Rebalance)
	assert.True(t, ethAsset.IsQuote)

	testAssetID, err := s.CreateAsset(
		"ABC",
		"ABC Advanced Token",
		ethereum.HexToAddress("0xB8c77482e45F1F44dE1745F52C74426C631bDD52"),
		12,
		commonv3.BTCFeed,
		true,
		false,
		testPWI,
		testRb,
		testAssetExchanges,
		testAssetTarget,
	)
	require.NoError(t, err)
	assert.Equal(t, uint64(2), testAssetID)

	assets, err = s.GetAssets()
	require.NoError(t, err)
	assert.Len(t, assets, 2)
	assert.NotEqual(t, assets[0].ID, assets[1].ID)
	testAsset := assets[1]
	assert.Equal(t, "ABC", testAsset.Symbol)
	assert.Equal(t, "ABC Advanced Token", testAsset.Name)
	assert.Equal(t, ethereum.HexToAddress("0xB8c77482e45F1F44dE1745F52C74426C631bDD52"), testAsset.Address)
	assert.Equal(t, uint64(12), testAsset.Decimals)
	assert.Equal(t, commonv3.BTCFeed, testAsset.SetRate)
	assert.True(t, testAsset.Rebalance)
	assert.False(t, testAsset.IsQuote)
	assert.Equal(t, testPWI, testAsset.PWI)
	assert.Equal(t, testRb, testAsset.RebalanceQuadratic)
	for i := range testAsset.Exchanges {
		assert.NotZero(t, testAsset.Exchanges[i].ID)
		// the input asset exchange does not have id field
		testAsset.Exchanges[i].ID = 0
	}
	assert.Equal(t, testAssetExchanges, testAsset.Exchanges)
	assert.Equal(t, testAssetTarget, testAsset.Target)
}

func TestStorage_GetAsset(t *testing.T) {
	db, tearDown := testutil.MustNewDevelopmentDB()
	defer func() {
		assert.NoError(t, tearDown())
	}()

	s, err := NewStorage(db)
	require.NoError(t, err)

	ethAsset, err := s.GetAsset(1)
	require.NoError(t, err)
	assert.Equal(t, uint64(1), ethAsset.ID)
	assert.Equal(t, "ETH", ethAsset.Symbol)
	assert.Equal(t, "Ethereum", ethAsset.Name)
	assert.Equal(t, ethereum.HexToAddress("0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"), ethAsset.Address)
	assert.Equal(t, uint64(18), ethAsset.Decimals)
	assert.Equal(t, commonv3.SetRateNotSet, ethAsset.SetRate)
	assert.False(t, ethAsset.Rebalance)
	assert.True(t, ethAsset.IsQuote)

	_, err = s.GetAsset(999)
	assert.Equal(t, commonv3.ErrNotFound, err)

	testAssetID, err := s.CreateAsset(
		"BTCXYZ",
		"Bitcoin Fork XYZ",
		ethereum.Address{},
		12,
		commonv3.BTCFeed,
		true,
		true,
		testPWI,
		testRb,
		testAssetExchanges,
		testAssetTarget,
	)
	require.NoError(t, err)

	testAsset, err := s.GetAsset(testAssetID)
	require.NoError(t, err)
	assert.Equal(t, "BTCXYZ", testAsset.Symbol)
	assert.Equal(t, "Bitcoin Fork XYZ", testAsset.Name)
	assert.Equal(t, ethereum.Address{}, testAsset.Address)
	assert.Equal(t, uint64(12), testAsset.Decimals)
	assert.Equal(t, commonv3.BTCFeed, testAsset.SetRate)
	assert.True(t, testAsset.Rebalance)
	assert.True(t, testAsset.IsQuote)
	assert.Equal(t, testPWI, testAsset.PWI)
	assert.Equal(t, testRb, testAsset.RebalanceQuadratic)
	for i := range testAsset.Exchanges {
		assert.NotZero(t, testAsset.Exchanges[i].ID)
		// the input asset exchange does not have id field
		testAsset.Exchanges[i].ID = 0
	}
	assert.Equal(t, testAssetExchanges, testAsset.Exchanges)
	assert.Equal(t, testAssetTarget, testAsset.Target)
}

func TestStorage_UpdateAsset(t *testing.T) {
	db, tearDown := testutil.MustNewDevelopmentDB()
	defer func() {
		assert.NoError(t, tearDown())
	}()

	s, err := NewStorage(db)
	require.NoError(t, err)

	testAssetID1, err := s.CreateAsset(
		"ABC",
		"ABC Advanced Token",
		ethereum.HexToAddress("0xB8c77482e45F1F44dE1745F52C74426C631bDD52"),
		12,
		commonv3.BTCFeed,
		true,
		false,
		testPWI,
		testRb,
		testAssetExchanges,
		testAssetTarget,
	)
	require.NoError(t, err)

	testAssetID2, err := s.CreateAsset(
		"DEF",
		"DEF Super Token",
		ethereum.HexToAddress("0xffe97fe10290715ba416a7c1Fd265F28dc574dd9"),
		12,
		commonv3.GoldFeed,
		true,
		true,
		testPWI,
		testRb,
		testAssetExchanges,
		testAssetTarget,
	)
	require.NoError(t, err)

	err = s.UpdateAsset(999, storage.WithSymbolUpdateAssetOption("random"))
	require.Equal(t, commonv3.ErrNotFound, err)

	err = s.UpdateAsset(testAssetID1, storage.WithSymbolUpdateAssetOption("DEF"))
	require.Equal(t, commonv3.ErrSymbolExists, err)

	testAsset1, err := s.GetAsset(testAssetID1)
	require.NoError(t, err)
	oldUpdated := testAsset1.Updated
	err = s.UpdateAsset(testAssetID1, storage.WithSymbolUpdateAssetOption("ABC2"))
	require.NoError(t, err)
	testAsset1, err = s.GetAsset(testAssetID1)
	require.NoError(t, err)
	assert.Equal(t, "ABC2", testAsset1.Symbol)
	assert.NotEqual(t, testAsset1.Updated, oldUpdated)

	// verify that we could have assets with same name
	err = s.UpdateAsset(testAssetID2, storage.WithNameUpdateAssetOption("ABC Advanced Token"))
	require.NoError(t, err)
	testAsset2, err := s.GetAsset(testAssetID2)
	require.NoError(t, err)
	assert.Equal(t, "ABC Advanced Token", testAsset2.Name)

	err = s.UpdateAsset(testAssetID2, storage.WithNameUpdateAssetOption("DEF Super Token 2"))
	require.NoError(t, err)
	testAsset2, err = s.GetAsset(testAssetID2)
	require.NoError(t, err)
	assert.Equal(t, "DEF Super Token 2", testAsset2.Name)

	err = s.UpdateAsset(
		testAssetID1,
		storage.WithAddressUpdateAssetOption(ethereum.HexToAddress("0xEA674fdDe714fd979de3EdF0F56AA9716B898ec8")))
	require.NoError(t, err)
	err = s.UpdateAsset(
		testAssetID2,
		storage.WithAddressUpdateAssetOption(ethereum.HexToAddress("0xEA674fdDe714fd979de3EdF0F56AA9716B898ec8")))
	assert.Equal(t, commonv3.ErrAddressExists, err)
	err = s.UpdateAsset(
		testAssetID2,
		storage.WithAddressUpdateAssetOption(ethereum.HexToAddress("0xea674fdde714fd979de3edf0f56aa9716b898ec8")))
	assert.Equal(t, commonv3.ErrAddressExists, err)

	err = s.UpdateAsset(
		testAssetID1,
		storage.WithDecimalsUpdateAssetOption(10),
	)
	require.NoError(t, err)
	testAsset1, err = s.GetAsset(testAssetID1)
	require.NoError(t, err)
	assert.Equal(t, uint64(10), testAsset1.Decimals)

	err = s.UpdateAsset(
		testAssetID1,
		storage.WithSetRateUpdateAssetOption(commonv3.SetRateNotSet),
	)
	require.NoError(t, err)
	testAsset1, err = s.GetAsset(testAssetID1)
	require.NoError(t, err)
	assert.Equal(t, commonv3.SetRateNotSet, testAsset1.SetRate)

	err = s.UpdateAsset(
		testAssetID1,
		storage.WithRebalanceUpdateAssetOption(false),
	)
	require.NoError(t, err)
	testAsset1, err = s.GetAsset(testAssetID1)
	require.NoError(t, err)
	assert.False(t, testAsset1.Rebalance)

	err = s.UpdateAsset(
		testAssetID1,
		storage.WithIsQuoteUpdateAssetOption(true),
	)
	require.NoError(t, err)
	testAsset1, err = s.GetAsset(testAssetID1)
	require.NoError(t, err)
	assert.True(t, testAsset1.IsQuote)

	testAsset1, err = s.GetAsset(testAssetID1)
	require.NoError(t, err)
	oldUpdated = testAsset1.Updated
	err = s.UpdateAsset(testAssetID1)
	require.NoError(t, err)
	assert.Equal(t, oldUpdated, testAsset1.Updated)
}

func TestStorage_ChangeAssetAddress(t *testing.T) {
	var (
		oldAddress = ethereum.HexToAddress("0xB8c77482e45F1F44dE1745F52C74426C631bDD52")
		newAddress = ethereum.HexToAddress("0xC2826E724Aa1cF01bC618B848453B2e0536F036E")
	)

	db, tearDown := testutil.MustNewDevelopmentDB()
	defer func() {
		assert.NoError(t, tearDown())
	}()

	s, err := NewStorage(db)
	require.NoError(t, err)

	testAssetID1, err := s.CreateAsset(
		"ABC",
		"ABC Advanced Token",
		oldAddress,
		12,
		commonv3.BTCFeed,
		true,
		false,
		testPWI,
		testRb,
		testAssetExchanges,
		testAssetTarget,
	)
	require.NoError(t, err)

	_, err = s.CreateAsset(
		"DEF",
		"DEF Super Token",
		ethereum.HexToAddress("0xffe97fe10290715ba416a7c1Fd265F28dc574dd9"),
		12,
		commonv3.GoldFeed,
		true,
		true,
		testPWI,
		testRb,
		testAssetExchanges,
		testAssetTarget,
	)
	require.NoError(t, err)

	err = s.ChangeAssetAddress(999, newAddress)
	require.Equal(t, commonv3.ErrNotFound, err)

	err = s.ChangeAssetAddress(testAssetID1, newAddress)
	require.NoError(t, err)
	testAsset1, err := s.GetAsset(testAssetID1)
	require.NoError(t, err)
	assert.Equal(t, newAddress, testAsset1.Address)
	assert.Equal(t, []ethereum.Address{oldAddress}, testAsset1.OldAddresses)
	assets, err := s.GetAssets()
	require.NoError(t, err)
	var found = false
	for _, assetDB := range assets {
		if assetDB.ID == testAssetID1 {
			testAsset1 = assetDB
			found = true
			break
		}
	}
	require.True(t, found)
	assert.Equal(t, []ethereum.Address{oldAddress}, testAsset1.OldAddresses)

	err = s.ChangeAssetAddress(testAssetID1, newAddress)
	require.Equal(t, commonv3.ErrAddressExists, err)

	err = s.ChangeAssetAddress(testAssetID1, ethereum.HexToAddress("0xffe97fe10290715ba416a7c1Fd265F28dc574dd9"))
	require.Equal(t, commonv3.ErrAddressExists, err)

	err = s.ChangeAssetAddress(testAssetID1, ethereum.HexToAddress("0xB8c77482e45F1F44dE1745F52C74426C631bDD52"))
	require.Equal(t, commonv3.ErrAddressExists, err)

	err = s.ChangeAssetAddress(testAssetID1, ethereum.HexToAddress("0x824AdA524aD4dd041036160F352a6F38411edF0B"))
	require.NoError(t, err)
	testAsset1, err = s.GetAsset(testAssetID1)
	require.NoError(t, err)
	assert.Equal(t, ethereum.HexToAddress("0x824AdA524aD4dd041036160F352a6F38411edF0B"), testAsset1.Address)
	assert.Equal(t, []ethereum.Address{oldAddress, newAddress}, testAsset1.OldAddresses)
	assets, err = s.GetAssets()
	require.NoError(t, err)
	found = false
	for _, assetDB := range assets {
		if assetDB.ID == testAssetID1 {
			testAsset1 = assetDB
			found = true
			break
		}
	}
	require.True(t, found)
	assert.Equal(t, []ethereum.Address{oldAddress, newAddress}, testAsset1.OldAddresses)
}

func TestStorage_GetTradingPairSymbols2(t *testing.T) {
	db, tearDown := testutil.MustNewDevelopmentDB()
	defer func() {
		assert.NoError(t, tearDown())
	}()

	s, err := NewStorage(db)
	require.NoError(t, err)

	nonQuoteAsset, err := s.CreateAsset(
		"NQT",
		"Non Quote Token",
		ethereum.HexToAddress("0x930e2f445A5a0e3c98b7C385125f95C24d772961"),
		12,
		commonv3.SetRateNotSet,
		false,
		false,
		nil,
		nil,
		nil,
		nil,
	)
	require.NoError(t, err)

	quoteAssetID, err := s.CreateAsset(
		"ABC",
		"ABC Advanced Token",
		ethereum.HexToAddress("0xB42b3F0C10385df057f7374d7Aa884571E71791b"),
		12,
		commonv3.SetRateNotSet,
		false,
		true,
		nil,
		nil,
		[]commonv3.AssetExchange{
			{
				ExchangeID:      uint64(common.Binance),
				Symbol:          "ABC",
				DepositAddress:  ethereum.HexToAddress("0x118ee757dD8841F81903E1C1d7d7Aa88e376cC39"),
				MinDeposit:      1,
				WithdrawFee:     2,
				PricePrecision:  3,
				AmountPrecision: 4,
				AmountLimitMin:  5,
				AmountLimitMax:  6,
				PriceLimitMin:   7,
				PriceLimitMax:   8,
				TradingPairs:    nil,
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
		commonv3.ExchangeFeed,
		true,
		false,
		testPWI,
		testRb,
		[]commonv3.AssetExchange{
			{
				ExchangeID:      uint64(common.Binance),
				Symbol:          "BNB",
				DepositAddress:  ethereum.HexToAddress("0x118ee757dD8841F81903E1C1d7d7Aa88e376cC39"),
				MinDeposit:      1,
				WithdrawFee:     2,
				PricePrecision:  3,
				AmountPrecision: 4,
				AmountLimitMin:  5,
				AmountLimitMax:  6,
				PriceLimitMin:   7,
				PriceLimitMax:   8,
				TradingPairs: []commonv3.TradingPair{
					{
						Base:        quoteAssetID,
						Quote:       quoteAssetID,
						MinNotional: 1.2,
					},
				},
			},
		},
		testAssetTarget,
	)
	require.Equal(t, commonv3.ErrBadTradingPairConfiguration, err)

	_, err = s.CreateAsset(
		"DEF",
		"DEF Super Token",
		ethereum.HexToAddress("0xffe97fe10290715ba416a7c1Fd265F28dc574dd9"),
		12,
		commonv3.ExchangeFeed,
		true,
		false,
		testPWI,
		testRb,
		[]commonv3.AssetExchange{
			{
				ExchangeID:      uint64(common.Binance),
				Symbol:          "BNB",
				DepositAddress:  ethereum.HexToAddress("0x118ee757dD8841F81903E1C1d7d7Aa88e376cC39"),
				MinDeposit:      1,
				WithdrawFee:     2,
				PricePrecision:  3,
				AmountPrecision: 4,
				AmountLimitMin:  5,
				AmountLimitMax:  6,
				PriceLimitMin:   7,
				PriceLimitMax:   8,
				TradingPairs: []commonv3.TradingPair{
					{
						Base:        0,
						Quote:       999, // not exist
						MinNotional: 1.3,
					},
				},
			},
		},
		testAssetTarget,
	)
	require.Equal(t, commonv3.ErrBadTradingPairConfiguration, err)

	_, err = s.CreateAsset(
		"DEF",
		"DEF Super Token",
		ethereum.HexToAddress("0xffe97fe10290715ba416a7c1Fd265F28dc574dd9"),
		12,
		commonv3.ExchangeFeed,
		true,
		false,
		testPWI,
		testRb,
		[]commonv3.AssetExchange{
			{
				ExchangeID:      uint64(common.Binance),
				Symbol:          "BNB",
				DepositAddress:  ethereum.HexToAddress("0x118ee757dD8841F81903E1C1d7d7Aa88e376cC39"),
				MinDeposit:      1,
				WithdrawFee:     2,
				PricePrecision:  3,
				AmountPrecision: 4,
				AmountLimitMin:  5,
				AmountLimitMax:  6,
				PriceLimitMin:   7,
				PriceLimitMax:   8,
				TradingPairs: []commonv3.TradingPair{
					{
						Base:        0,
						Quote:       nonQuoteAsset,
						MinNotional: 1.4,
					},
				},
			},
		},
		testAssetTarget,
	)
	require.Equal(t, commonv3.ErrBadTradingPairConfiguration, err)

	testAssetID, err := s.CreateAsset(
		"DEF",
		"DEF Super Token",
		ethereum.HexToAddress("0xffe97fe10290715ba416a7c1Fd265F28dc574dd9"),
		12,
		commonv3.ExchangeFeed,
		true,
		false,
		testPWI,
		testRb,
		[]commonv3.AssetExchange{
			{
				ExchangeID:      uint64(common.Binance),
				Symbol:          "BNB",
				DepositAddress:  ethereum.HexToAddress("0x118ee757dD8841F81903E1C1d7d7Aa88e376cC39"),
				MinDeposit:      1,
				WithdrawFee:     2,
				PricePrecision:  3,
				AmountPrecision: 4,
				AmountLimitMin:  5,
				AmountLimitMax:  6,
				PriceLimitMin:   7,
				PriceLimitMax:   8,
				TradingPairs: []commonv3.TradingPair{
					{
						Base:        0,
						Quote:       quoteAssetID,
						MinNotional: 1.5,
					},
				},
			},
		},
		testAssetTarget,
	)
	require.NoError(t, err)

	testAsset, err := s.GetAsset(testAssetID)
	require.NoError(t, err)
	require.Len(t, testAsset.Exchanges, 1)
	require.NotNil(t, testAsset.Exchanges[0].TradingPairs)
	require.Len(t, testAsset.Exchanges[0].TradingPairs, 1)
	assert.Equal(t, testAssetID, testAsset.Exchanges[0].TradingPairs[0].Base)
	assert.Equal(t, quoteAssetID, testAsset.Exchanges[0].TradingPairs[0].Quote)
	assert.Equal(t, 1.5, testAsset.Exchanges[0].TradingPairs[0].MinNotional)

	assets, err := s.GetAssets()
	require.NoError(t, err)
	for _, asset := range assets {
		if asset.ID == testAssetID {
			testAsset = asset
			require.Len(t, testAsset.Exchanges, 1)
			require.NotNil(t, testAsset.Exchanges[0].TradingPairs)
			require.Len(t, testAsset.Exchanges[0].TradingPairs, 1)
			assert.Equal(t, testAssetID, testAsset.Exchanges[0].TradingPairs[0].Base)
			assert.Equal(t, quoteAssetID, testAsset.Exchanges[0].TradingPairs[0].Quote)
			assert.Equal(t, 1.5, testAsset.Exchanges[0].TradingPairs[0].MinNotional)
		}
	}

	pairs, err := s.GetTradingPairSymbols(uint64(common.Binance))
	require.NoError(t, err)
	require.Len(t, pairs, 1)
	assert.Equal(t, "BNB", pairs[0].BaseSymbol)
	assert.Equal(t, "ABC", pairs[0].QuoteSymbol)
}

func TestStorage_GetMinNotional(t *testing.T) {
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
		commonv3.SetRateNotSet,
		false,
		true,
		nil,
		nil,
		[]commonv3.AssetExchange{
			{
				ExchangeID:      uint64(common.Binance),
				Symbol:          "ABC",
				DepositAddress:  ethereum.HexToAddress("0x118ee757dD8841F81903E1C1d7d7Aa88e376cC39"),
				MinDeposit:      1,
				WithdrawFee:     2,
				PricePrecision:  3,
				AmountPrecision: 4,
				AmountLimitMin:  5,
				AmountLimitMax:  6,
				PriceLimitMin:   7,
				PriceLimitMax:   8,
				TradingPairs:    nil,
			},
		},
		testAssetTarget,
	)
	require.NoError(t, err)

	testAssetID, err := s.CreateAsset(
		"DEF",
		"DEF Super Token",
		ethereum.HexToAddress("0xffe97fe10290715ba416a7c1Fd265F28dc574dd9"),
		12,
		commonv3.ExchangeFeed,
		true,
		false,
		testPWI,
		testRb,
		[]commonv3.AssetExchange{
			{
				ExchangeID:      uint64(common.Binance),
				Symbol:          "BNB",
				DepositAddress:  ethereum.HexToAddress("0x118ee757dD8841F81903E1C1d7d7Aa88e376cC39"),
				MinDeposit:      1,
				WithdrawFee:     2,
				PricePrecision:  3,
				AmountPrecision: 4,
				AmountLimitMin:  5,
				AmountLimitMax:  6,
				PriceLimitMin:   7,
				PriceLimitMax:   8,
				TradingPairs: []commonv3.TradingPair{
					{
						Base:        0,
						Quote:       quoteAssetID,
						MinNotional: 1.5,
					},
				},
			},
		},
		testAssetTarget,
	)
	require.NoError(t, err)

	minNotional, err := s.GetMinNotional(uint64(common.Binance), testAssetID, quoteAssetID)
	require.NoError(t, err)
	assert.Equal(t, 1.5, minNotional)
	_, err = s.GetMinNotional(uint64(common.Binance), quoteAssetID, 1)
	require.Equal(t, commonv3.ErrNotFound, err)
}

func TestStorage_Initialization(t *testing.T) {
	db, tearDown := testutil.MustNewDevelopmentDB()
	defer func() {
		assert.NoError(t, tearDown())
	}()

	// make sure the database can be initialization multiple times
	_, err := NewStorage(db)
	require.NoError(t, err)
	_, err = NewStorage(db)
	require.NoError(t, err)
}
