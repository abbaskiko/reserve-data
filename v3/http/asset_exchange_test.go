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

func TestHTTPServerAssetExchange(t *testing.T) {

	const (
		createAssetExchange = "/v3/create-asset-exchange"
		// updateAssetExchange = "/v3/update-asset-exchange"
		// updateTradingPair = "/v3/update-trading-pair"
	)

	db, tearDown := testutil.MustNewDevelopmentDB()
	defer func() {
		assert.NoError(t, tearDown())
	}()

	s, err := postgres.NewStorage(db)
	require.NoError(t, err)

	server := NewServer(s, nil)

	var tests = []testCase{

		{
			msg:      "create asset exchange",
			endpoint: createAssetExchange,
			method:   http.MethodPost,
			assert:   httputil.ExpectSuccess,
			data: &common.CreateCreateAssetExchange{
				AssetExchanges: []common.CreateAssetExchangeEntry{
					{
						AssetID:           1,
						ExchangeID:        0,
						Symbol:            "ETH",
						DepositAddress:    eth.HexToAddress("0x00"),
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
		// TODO add more test here
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.msg, func(t *testing.T) { testHTTPRequest(t, tc, server.r) })
	}
}
