package postgres

import (
	"testing"

	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KyberNetwork/reserve-data/common/testutil"
	commonv3 "github.com/KyberNetwork/reserve-data/v3/common"
	"github.com/KyberNetwork/reserve-data/v3/storage"
)

func TestStorage_CreateAssetExchange(t *testing.T) {
	db, tearDown := testutil.MustNewDevelopmentDB()
	defer func() {
		assert.NoError(t, tearDown())
	}()

	s, err := NewStorage(db)
	require.NoError(t, err)

	var tests = []struct {
		msg string

		exchangeID     uint64
		assetID        uint64
		symbol         string
		depositAddress ethereum.Address

		minDeposit        float64
		withdrawFee       float64
		targetRecommended float64
		targetRatio       float64

		assertFn func(*testing.T, uint64, error)
	}{
		{
			msg:               "create asset exchange with valid data",
			exchangeID:        1,
			assetID:           1,
			symbol:            "ETH",
			depositAddress:    ethereum.Address{},
			minDeposit:        0.0001,
			withdrawFee:       0.0001,
			targetRecommended: 0.001,
			targetRatio:       0.001,

			assertFn: func(t *testing.T, id uint64, err error) {
				assert.NotZero(t, id)
				require.Equal(t, nil, err)
			},
		},
		{
			msg:               "create asset exchange with invalid constraints",
			exchangeID:        1,
			assetID:           1,
			symbol:            "ETH",
			depositAddress:    ethereum.Address{},
			minDeposit:        0.0001,
			withdrawFee:       0.0001,
			targetRecommended: 0.001,
			targetRatio:       0.001,

			assertFn: func(t *testing.T, id uint64, err error) {
				assert.Zero(t, id)
				require.Equal(t, commonv3.ErrDuplicateExchangeIDAssetID, err)
			},
		},
		{
			msg:               "create asset exchange with not existed exchange id",
			exchangeID:        1234,
			assetID:           1,
			symbol:            "ETH",
			depositAddress:    ethereum.Address{},
			minDeposit:        0.0001,
			withdrawFee:       0.0001,
			targetRecommended: 0.001,
			targetRatio:       0.001,

			assertFn: func(t *testing.T, id uint64, err error) {
				assert.Zero(t, id)
				require.Equal(t, commonv3.ErrExchangeIDNotExists, err)
			},
		},
		{
			msg:               "create asset exchange with not existed asset id",
			exchangeID:        1,
			assetID:           1234,
			symbol:            "ETH",
			depositAddress:    ethereum.Address{},
			minDeposit:        0.0001,
			withdrawFee:       0.0001,
			targetRecommended: 0.001,
			targetRatio:       0.001,

			assertFn: func(t *testing.T, id uint64, err error) {
				assert.Zero(t, id)
				require.Equal(t, commonv3.ErrAssetIDNotExists, err)
			},
		},
	}

	for _, tc := range tests {
		t.Logf("running test case for: %s", tc.msg)
		id, err := s.CreateAssetExchange(
			tc.exchangeID,
			tc.assetID,
			tc.symbol,
			tc.depositAddress,
			tc.minDeposit,
			tc.withdrawFee,
			tc.targetRecommended,
			tc.targetRatio,
		)
		tc.assertFn(t, id, err)
	}
}

func TestStorage_UpdateAssetExchange(t *testing.T) {
	db, tearDown := testutil.MustNewDevelopmentDB()
	defer func() {
		assert.NoError(t, tearDown())
	}()

	s, err := NewStorage(db)
	require.NoError(t, err)

	assetExchangeID, err := s.CreateAssetExchange(1, 1, "ETH", ethereum.Address{}, 0.0001, 0.0001, 0.001, 0.001)
	require.NoError(t, err)
	require.NotZero(t, assetExchangeID)

	var tests = []struct {
		msg  string
		id   uint64
		opts storage.UpdateAssetExchangeOpts

		assertFn func(*testing.T, error)
	}{
		{
			msg: "test update not existed id",
			id:  123,
			opts: storage.UpdateAssetExchangeOpts{
				Symbol: commonv3.StringPointer("KNC"),
			},
			assertFn: func(t *testing.T, err error) {
				require.Equal(t, commonv3.ErrNotFound, err)
			},
		},
		{
			msg: "test update symbol",
			id:  assetExchangeID,
			opts: storage.UpdateAssetExchangeOpts{
				Symbol: commonv3.StringPointer("KNC"),
			},
			assertFn: func(t *testing.T, err error) {
				updatedAE, errG := s.GetAssetExchange(assetExchangeID)
				require.NoError(t, errG)
				require.Equal(t, updatedAE.Symbol, "KNC")
			},
		},
		{
			msg: "test update deposit address",
			id:  assetExchangeID,
			opts: storage.UpdateAssetExchangeOpts{
				DepositAddress: commonv3.AddressPointer(ethereum.HexToAddress("0xea674fdde714fd979de3edf0f56aa9716b898ec8")),
			},
			assertFn: func(t *testing.T, err error) {
				updatedAE, errG := s.GetAssetExchange(assetExchangeID)
				require.NoError(t, errG)
				require.Equal(t, updatedAE.DepositAddress, ethereum.HexToAddress("0xea674fdde714fd979de3edf0f56aa9716b898ec8"))
			},
		},
		{
			msg: "test update min deposit",
			id:  assetExchangeID,
			opts: storage.UpdateAssetExchangeOpts{
				MinDeposit: commonv3.FloatPointer(0.0002),
			},
			assertFn: func(t *testing.T, err error) {
				updatedAE, errG := s.GetAssetExchange(assetExchangeID)
				require.NoError(t, errG)
				require.Equal(t, updatedAE.MinDeposit, 0.0002)
			},
		},
		{
			msg: "test update withdraw fee",
			id:  assetExchangeID,
			opts: storage.UpdateAssetExchangeOpts{
				WithdrawFee: commonv3.FloatPointer(0.0002),
			},
			assertFn: func(t *testing.T, err error) {
				updatedAE, errG := s.GetAssetExchange(assetExchangeID)
				require.NoError(t, errG)
				require.Equal(t, updatedAE.WithdrawFee, 0.0002)
			},
		},
		{
			msg: "test update target recommended",
			id:  assetExchangeID,
			opts: storage.UpdateAssetExchangeOpts{
				TargetRecommended: commonv3.FloatPointer(0.002),
			},
			assertFn: func(t *testing.T, err error) {
				updatedAE, errG := s.GetAssetExchange(assetExchangeID)
				require.NoError(t, errG)
				require.Equal(t, updatedAE.TargetRecommended, 0.002)
			},
		},
		{
			msg: "test update target ratio",
			id:  assetExchangeID,
			opts: storage.UpdateAssetExchangeOpts{
				TargetRatio: commonv3.FloatPointer(0.002),
			},
			assertFn: func(t *testing.T, err error) {
				updatedAE, errG := s.GetAssetExchange(assetExchangeID)
				require.NoError(t, errG)
				require.Equal(t, updatedAE.TargetRatio, 0.002)
			},
		},
	}

	for _, tc := range tests {
		t.Logf("running test case for: %s", tc.msg)
		err = s.UpdateAssetExchange(tc.id, tc.opts)
		tc.assertFn(t, err)
	}
}
