package http

import (
	"net/http"
	"testing"

	eth "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KyberNetwork/reserve-data/common/testutil"
	"github.com/KyberNetwork/reserve-data/http/httputil"
	"github.com/KyberNetwork/reserve-data/v3/common"
	"github.com/KyberNetwork/reserve-data/v3/storage/postgres"
)

func TestHTTPServerTradingPair(t *testing.T) {

	const (
		assetBase         = "/v3/asset"
		createAssetBase   = "/v3/create-asset"
		createTradingPair = "/v3/create-trading-pair"
		// updateTradingPair = "/v3/update-trading-pair"
	)

	db, tearDown := testutil.MustNewDevelopmentDB()
	defer func() {
		assert.NoError(t, tearDown())
	}()

	s, err := postgres.NewStorage(db)
	require.NoError(t, err)
	// asset = 1 for ETH is pre-insert in DB.
	_, err = s.CreateAssetExchange(0, 1, "ETH", eth.HexToAddress("0x00"), 10,
		0.2, 5.0, 0.3)
	require.NoError(t, err)
	server := NewServer(s, nil)

	var tests = []testCase{

		{
			msg:      "create pending asset",
			endpoint: createAssetBase,
			method:   http.MethodPost,
			assert:   httputil.ExpectSuccess,
			data:     createPEA,
		},
		{
			msg:      "confirm pending asset",
			endpoint: createAssetBase + "/1",
			method:   http.MethodPut,
			data:     nil,
			assert:   httputil.ExpectSuccess,
		},
		{
			msg:      "receive asset",
			endpoint: assetBase + "/2",
			method:   http.MethodGet,
			assert:   httputil.ExpectSuccess,
		},
		{
			msg:      "create trading pair",
			endpoint: createTradingPair,
			method:   http.MethodPost,
			data: common.CreateCreateTradingPair{
				TradingPairs: []common.CreateTradingPairEntry{
					{
						TradingPair: common.TradingPair{
							Base:            0,
							Quote:           1,
							PricePrecision:  1.0,
							AmountPrecision: 1.0,
							AmountLimitMin:  100.0,
							AmountLimitMax:  1000.0,
							PriceLimitMin:   100.0,
							PriceLimitMax:   1000.0,
							MinNotional:     100.0,
						},
						ExchangeID: 1,
					},
				},
			},
			assert: httputil.ExpectSuccess,
		},
		{
			msg:      "confirm create trading pair",
			endpoint: createTradingPair + "/1",
			method:   http.MethodPut,
			assert:   httputil.ExpectFailure, // will fail because we did not config assets on exchange id=1
		},
		// TODO add more test here
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.msg, func(t *testing.T) { testHTTPRequest(t, tc, server.r) })
	}
}
