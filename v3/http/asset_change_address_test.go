package http

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/KyberNetwork/reserve-data/common/testutil"
	"github.com/KyberNetwork/reserve-data/http/httputil"
	"github.com/KyberNetwork/reserve-data/v3/common"

	"github.com/KyberNetwork/reserve-data/v3/storage/postgres"
)

func TestHTTPServerChangeAssetAddress(t *testing.T) {

	t.Skip()
	const (
		changeAssetAddress = "/v3/change-asset-address"
	)

	db, tearDown := testutil.MustNewDevelopmentDB()
	defer func() {
		assert.NoError(t, tearDown())
	}()

	s, err := postgres.NewStorage(db)
	require.NoError(t, err)

	assetID, err := createSampleAsset(s)
	require.NoError(t, err)

	asset, err := s.GetAsset(assetID)
	require.NoError(t, err)

	assetAddress := asset.Address.Hex()
	t.Logf("address already create in change address %s", assetAddress)

	server := NewServer(s, nil)

	var tests = []testCase{
		{
			msg:      "create change asset address",
			endpoint: changeAssetAddress,
			method:   http.MethodPost,
			assert:   httputil.ExpectSuccess,
			data: &common.ChangeAssetAddress{
				Assets: []common.ChangeAssetAddressEntry{
					{
						ID:      assetID,
						Address: "0x4ddbda50ddbec0289a80de667a72d158819e381d",
					},
				},
			},
		},
		{
			msg:      "confirm pending change asset address",
			endpoint: changeAssetAddress + "/1",
			method:   http.MethodPut,
			data:     nil,
			assert:   httputil.ExpectSuccess,
		},
		{
			msg: "create invalid change asset address (with invalid address)",
			data: &common.ChangeAssetAddress{
				Assets: []common.ChangeAssetAddressEntry{
					{
						ID:      assetID,
						Address: "invalid address",
					},
				},
			},
			endpoint: changeAssetAddress,
			method:   http.MethodPost,
			assert:   httputil.ExpectFailureWithReason(common.ErrInvalidAddress.Error() + "fsdf"),
		},
		{
			msg: "create invalid change asset address (with same current address)",
			data: &common.ChangeAssetAddress{
				Assets: []common.ChangeAssetAddressEntry{
					{
						ID:      assetID,
						Address: assetAddress,
					},
				},
			},
			endpoint: changeAssetAddress,
			method:   http.MethodPost,
			assert:   httputil.ExpectFailure,
		},
		{
			msg:      "create change other asset address",
			endpoint: changeAssetAddress,
			method:   http.MethodPost,
			assert:   httputil.ExpectSuccess,
			data: &common.ChangeAssetAddress{
				Assets: []common.ChangeAssetAddressEntry{
					{
						ID:      assetID,
						Address: "0x5dbcb95364cbc5604bacbb8c6eb9aa788f347a17",
					},
				},
			},
		},
		{
			msg:      "reject pending change asset address",
			endpoint: changeAssetAddress + "/2",
			method:   http.MethodDelete,
			data:     nil,
			assert:   httputil.ExpectSuccess,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.msg, func(t *testing.T) { testHTTPRequest(t, tc, server.r) })
	}
}
