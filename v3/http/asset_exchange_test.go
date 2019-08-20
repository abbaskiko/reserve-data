package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	eth "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	v1common "github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/common/testutil"
	"github.com/KyberNetwork/reserve-data/http/httputil"
	"github.com/KyberNetwork/reserve-data/v3/common"
	"github.com/KyberNetwork/reserve-data/v3/storage/postgres"
)

func TestHTTPServerAssetExchange(t *testing.T) {

	const (
		createAssetExchange = "/v3/create-asset-exchange"
		updateAssetExchange = "/v3/update-asset-exchange"
		// updateTradingPair = "/v3/update-trading-pair"
	)

	db, tearDown := testutil.MustNewDevelopmentDB()
	defer func() {
		assert.NoError(t, tearDown())
	}()

	s, err := postgres.NewStorage(db)
	require.NoError(t, err)

	assetID, err := createSampleAsset(s)
	require.NoError(t, err)

	server := NewServer(s, "", nil)

	var tests = []testCase{

		{
			msg:      "create asset exchange",
			endpoint: createAssetExchange,
			method:   http.MethodPost,
			assert:   httputil.ExpectSuccess,
			data: &common.CreateCreateAssetExchange{
				AssetExchanges: []common.CreateAssetExchangeEntry{
					{
						AssetID:           assetID,
						ExchangeID:        huobi,
						Symbol:            "ETH",
						DepositAddress:    eth.HexToAddress("0x001"),
						MinDeposit:        10.0,
						WithdrawFee:       11.0,
						TargetRecommended: 12.0,
						TargetRatio:       13.0,
					},
				}},
		},
		{
			msg:      "confirm pending asset exchange",
			endpoint: createAssetExchange + "/1",
			method:   http.MethodPut,
			data:     nil,
			assert:   httputil.ExpectSuccess,
		},
		{ // test deposit addr = 0
			data: common.CreateUpdateAssetExchange{
				AssetExchanges: []common.UpdateAssetExchangeEntry{
					{
						ID:             1,
						DepositAddress: common.AddressPointer(eth.HexToAddress("0x0000000000000")),
						MinDeposit:     common.FloatPointer(6.0),
						WithdrawFee:    common.FloatPointer(9.0),
					},
				},
			},
			endpoint: updateAssetExchange,
			method:   http.MethodPost,
			msg:      "create invalid update asset exchange",
			assert:   httputil.ExpectFailureWithReason(common.ErrDepositAddressMissing.Error()),
		},
		{ // test create update asset exchange
			data: common.CreateUpdateAssetExchange{
				AssetExchanges: []common.UpdateAssetExchangeEntry{
					{
						ID:             2,
						DepositAddress: common.AddressPointer(eth.HexToAddress("0x0000000000001")),
						MinDeposit:     common.FloatPointer(6.0),
						WithdrawFee:    common.FloatPointer(9.0),
					},
				},
			},
			endpoint: updateAssetExchange,
			method:   http.MethodPost,
			msg:      "create update asset exchange",
			assert:   httputil.ExpectSuccess,
		},
		{ // test confirm update asset exchange
			msg:      "confirm pending update asset exchange",
			endpoint: updateAssetExchange + "/2",
			method:   http.MethodPut,
			data:     nil,
			assert:   httputil.ExpectSuccess,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.msg, func(t *testing.T) { testHTTPRequest(t, tc, server.r) })
	}
}

func TestHTTPServer_DeleteAssetExchange(t *testing.T) {
	const (
		settingChangePath = "/v3/setting-change"
	)

	var (
		supportedExchanges = make(map[v1common.ExchangeID]v1common.LiveExchange)
	)

	//create map of test exchange
	for _, exchangeID := range []v1common.ExchangeID{v1common.Binance, v1common.Huobi, v1common.StableExchange} {
		exchange := v1common.TestExchange{}
		supportedExchanges[exchangeID] = exchange
	}

	db, tearDown := testutil.MustNewDevelopmentDB()
	defer func() {
		assert.NoError(t, tearDown())
	}()

	s, err := postgres.NewStorage(db)
	require.NoError(t, err)

	assetID, err := createSampleAsset(s)
	require.NoError(t, err)

	server := NewServer(s, "", supportedExchanges)

	//create map of test exchange
	for _, exchangeID := range []v1common.ExchangeID{v1common.Binance, v1common.Huobi, v1common.StableExchange} {
		exchange := v1common.TestExchange{}
		v1common.SupportedExchanges[exchangeID] = exchange
	}

	var tests = []testCase{

		{
			msg:      "create asset exchange",
			endpoint: settingChangePath,
			method:   http.MethodPost,
			assert:   httputil.ExpectSuccess,
			data: &common.SettingChange{
				ChangeList: []common.SettingChangeEntry{
					{
						Type: common.ChangeTypeCreateAssetExchange,
						Data: common.CreateAssetExchangeEntry{
							AssetID:           assetID,
							ExchangeID:        huobi,
							Symbol:            "ETH",
							DepositAddress:    eth.HexToAddress("0x001"),
							MinDeposit:        10.0,
							WithdrawFee:       11.0,
							TargetRecommended: 12.0,
							TargetRatio:       13.0,
						},
					},
				},
			},
		},
		{
			msg:      "confirm pending asset exchange",
			endpoint: settingChangePath + "/1",
			method:   http.MethodPut,
			data:     nil,
			assert:   httputil.ExpectSuccess,
		},
		{ // test deposit addr = 0
			data: common.SettingChange{
				ChangeList: []common.SettingChangeEntry{
					{
						Type: common.ChangeTypeUpdateAssetExchange,
						Data: common.UpdateAssetExchangeEntry{
							ID:             1,
							DepositAddress: common.AddressPointer(eth.HexToAddress("0x0000000000000")),
							MinDeposit:     common.FloatPointer(6.0),
							WithdrawFee:    common.FloatPointer(9.0),
						},
					},
				},
			},
			endpoint: settingChangePath,
			method:   http.MethodPost,
			msg:      "create invalid update asset exchange",
			assert:   httputil.ExpectFailure,
		},
		{ // test create update asset exchange
			data: common.SettingChange{
				ChangeList: []common.SettingChangeEntry{
					{
						Type: common.ChangeTypeUpdateAssetExchange,
						Data: common.UpdateAssetExchangeEntry{
							ID:             2,
							DepositAddress: common.AddressPointer(eth.HexToAddress("0x0000000000001")),
							MinDeposit:     common.FloatPointer(6.0),
							WithdrawFee:    common.FloatPointer(9.0),
						},
					},
				},
			},
			endpoint: settingChangePath,
			method:   http.MethodPost,
			msg:      "create update asset exchange",
			assert:   httputil.ExpectSuccess,
		},
		{ // test confirm update asset exchange
			msg:      "confirm pending update asset exchange",
			endpoint: settingChangePath + "/2",
			method:   http.MethodPut,
			data:     nil,
			assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				httputil.ExpectSuccess(t, resp)
				ae, err := s.GetAssetExchange(2)
				assert.NoError(t, err)
				assert.Equal(t, float64(6.0), ae.MinDeposit)
				fmt.Printf("len(trading pairs)=%d\n", len(ae.TradingPairs))
			},
		},
		{
			msg:      "create invalid delete asset exchange request",
			endpoint: settingChangePath,
			data: common.SettingChange{
				ChangeList: []common.SettingChangeEntry{
					{
						Type: common.ChangeTypeDeleteAssetExchange,
						Data: common.DeleteAssetExchangeEntry{
							AssetExchangeID: 2,
						},
					},
				},
			},
			method: http.MethodPost,
			assert: httputil.ExpectFailure,
		},
		{
			msg:      "confirm invalid delete asset exchange request",
			endpoint: settingChangePath + "/3",
			data:     nil,
			method:   http.MethodPut,
			assert:   httputil.ExpectFailure,
		},
		{
			msg:      "create valid delete asset exchange request",
			endpoint: settingChangePath,
			data: common.SettingChange{
				ChangeList: []common.SettingChangeEntry{
					{
						Type: common.ChangeTypeDeleteAssetExchange,
						Data: common.DeleteAssetExchangeEntry{
							AssetExchangeID: 4,
						},
					},
				},
			},
			method: http.MethodPost,
			assert: httputil.ExpectSuccess,
		},
		{
			msg:      "confirm valid delete asset exchange request",
			endpoint: settingChangePath + "/3",
			data:     nil,
			method:   http.MethodPut,
			assert:   httputil.ExpectSuccess,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.msg, func(t *testing.T) { testHTTPRequest(t, tc, server.r) })
	}
}
