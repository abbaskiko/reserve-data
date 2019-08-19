package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KyberNetwork/reserve-data/common/testutil"
	"github.com/KyberNetwork/reserve-data/http/httputil"
	"github.com/KyberNetwork/reserve-data/v3/common"
	"github.com/KyberNetwork/reserve-data/v3/storage/postgres"
)

func TestPriceFactor(t *testing.T) {
	db, tearDown := testutil.MustNewDevelopmentDB()
	defer func() {
		assert.NoError(t, tearDown())
	}()
	s, err := postgres.NewStorage(db)
	require.NoError(t, err)

	server := NewServer(s, "")
	const priceFactor = "/v3/price-factor"
	var tests = []testCase{
		{
			msg:      "create price factor 1",
			endpoint: priceFactor,
			method:   http.MethodPost,
			data: common.PriceFactorAtTime{
				Timestamp: 3,
				Data: []common.AssetPriceFactor{
					{
						AssetID: 1,
						AfpMid:  31,
						Spread:  32,
					},
					{
						AssetID: 2,
						AfpMid:  33,
						Spread:  34,
					},
				},
			},
			assert: httputil.ExpectSuccess,
		},
		{
			msg:      "create price factor 2",
			endpoint: priceFactor,
			method:   http.MethodPost,
			data: common.PriceFactorAtTime{
				Timestamp: 4,
				Data: []common.AssetPriceFactor{
					{
						AssetID: 1,
						AfpMid:  31,
						Spread:  32,
					},
					{
						AssetID: 2,
						AfpMid:  33,
						Spread:  34,
					},
				},
			},
			assert: httputil.ExpectSuccess,
		},
		{
			msg:      "create price factor 3",
			endpoint: priceFactor,
			method:   http.MethodPost,
			data: common.PriceFactorAtTime{
				Timestamp: 5,
				Data: []common.AssetPriceFactor{
					{
						AssetID: 1,
						AfpMid:  31,
						Spread:  32,
					},
					{
						AssetID: 2,
						AfpMid:  33,
						Spread:  34,
					},
				},
			},
			assert: httputil.ExpectSuccess,
		},
		{
			msg:      "select only 2",
			endpoint: priceFactor + "?from=3&to=4",
			method:   http.MethodGet,
			assert: func(t *testing.T, resp *httptest.ResponseRecorder) {
				var res struct {
					Data []common.AssetPriceFactorListResponse `json:"data"`
				}
				assert.Equal(t, resp.Code, http.StatusOK)
				err := json.NewDecoder(resp.Body).Decode(&res)
				assert.NoError(t, err)
				assert.Len(t, res.Data, 2)         // expect 2 asset in response
				assert.Len(t, res.Data[0].Data, 2) // expect 2 values for one asset.
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.msg, func(t *testing.T) { testHTTPRequest(t, tc, server.r) })
	}
}
