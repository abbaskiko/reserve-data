package http

import (
	"encoding/json"
	"fmt"
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

var (
	createPEA = common.CreatePendingAsset{AssetInputs: []common.CreatePendingAssetEntry{
		{
			Address:   eth.HexToAddress("0x01"),
			Symbol:    "ABC",
			Name:      "ABC test",
			Decimals:  18,
			SetRate:   common.ExchangeFeed,
			Rebalance: true,
			IsQuote:   true,
			PWI: &common.AssetPWI{
				Ask: common.PWIEquation{},
				Bid: common.PWIEquation{},
			},
			RebalanceQuadratic: &common.RebalanceQuadratic{},
			Target:             &common.AssetTarget{},
			Exchanges: []common.AssetExchange{
				{
					TradingPairs: []common.TradingPair{
						{
							Base:  1,
							Quote: 0,
						},
					},
				},
			},
		},
	}}
)

func TestReCreatePendingAsset(t *testing.T) {

	db, tearDown := testutil.MustNewDevelopmentDB()
	defer func() {
		assert.NoError(t, tearDown())
	}()
	s, err := postgres.NewStorage(db)
	require.NoError(t, err)
	jsonPayload, err := json.Marshal(createPEA)
	require.NoError(t, err)
	_, err = s.CreatePendingAsset(jsonPayload)
	require.NoError(t, err)
	id2, err := s.CreatePendingAsset(jsonPayload)
	require.NoError(t, err)
	pending, err := s.ListPendingAsset()
	require.NoError(t, err)
	if len(pending) != 1 || pending[0].ID != id2 {
		t.Fatal("expect 1 element with latest create one")
	}
}

func TestHTTPServerAsset(t *testing.T) {

	const (
		assetBase        = "/v3/asset"
		pendingAssetBase = "/v3/pending-asset"
	)

	db, tearDown := testutil.MustNewDevelopmentDB()
	defer func() {
		assert.NoError(t, tearDown())
	}()

	s, err := postgres.NewStorage(db)
	require.NoError(t, err)
	assetID := 1
	_, err = s.CreateAssetExchange(0, 1, "ETH", eth.HexToAddress("0x00"), 10,
		0.2, 5.0, 0.3)
	require.NoError(t, err)
	server := NewServer(s, nil)

	var tests = []testCase{
		{
			msg:      "asset not found",
			endpoint: assetBase + "/-1",
			method:   http.MethodGet,
			assert:   httputil.ExpectFailure,
		},
		{
			msg:      "receive asset",
			endpoint: assetBase + fmt.Sprintf("/%d", assetID),
			method:   http.MethodGet,
			assert:   httputil.ExpectFailure,
		},
		{
			msg:      "list all asset",
			endpoint: assetBase,
			method:   http.MethodGet,
			assert:   httputil.ExpectSuccess,
		},
		{
			msg:      "create pending asset",
			endpoint: pendingAssetBase,
			method:   http.MethodPost,
			assert:   newAssertCreated(nil),
			data:     createPEA,
		},
		{
			msg:      "list all pending asset",
			endpoint: pendingAssetBase,
			method:   http.MethodGet,
			assert:   newAssertHTTPCode(http.StatusOK),
		},
		{
			msg:      "confirm invalid pending asset",
			endpoint: pendingAssetBase + "/-1",
			method:   http.MethodPut,
			data:     nil,
			assert:   newAssertHTTPCode(http.StatusBadRequest),
		},
		{
			msg:      "confirm pending asset",
			endpoint: pendingAssetBase + "/1",
			method:   http.MethodPut,
			data:     nil,
			assert:   newAssertHTTPCode(http.StatusOK),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.msg, func(t *testing.T) { testHTTPRequest(t, tc, server.r) })
	}
}
