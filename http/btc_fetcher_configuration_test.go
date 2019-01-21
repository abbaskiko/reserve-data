package http

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/KyberNetwork/reserve-data/core"
	"github.com/KyberNetwork/reserve-data/data"
	"github.com/KyberNetwork/reserve-data/data/storage"
	"github.com/KyberNetwork/reserve-data/http/httputil"
	"github.com/KyberNetwork/reserve-data/settings"
	"github.com/gin-gonic/gin"
)

func TestHTTPServerBTCConfiguration(t *testing.T) {
	t.Skip()
	const (
		updateBTCFetcherConfigurationEndpoint = "/set-btc-fetcher-configuration"
		getBTCFetcherConfigurationEndpoint    = "/get-btc-fetcher-configuration"
	)
	tmpDir, err := ioutil.TempDir("", "test_btc_fetcher_configuration")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if rErr := os.RemoveAll(tmpDir); rErr != nil {
			t.Error(rErr)
		}
	}()

	testStorage, err := storage.NewBoltStorage(filepath.Join(tmpDir, "setting.db"))
	if err != nil {
		log.Fatal(err)
	}
	tokenSetting, err := settings.NewTokenSetting(nil)
	if err != nil {
		log.Fatal(err)
	}
	addressSetting := &settings.AddressSetting{}

	exchangeSetting, err := settings.NewExchangeSetting(nil)
	if err != nil {
		log.Fatal(err)
	}

	setting, err := settings.NewSetting(tokenSetting, addressSetting, exchangeSetting)
	if err != nil {
		log.Fatal(err)
	}

	testServer := HTTPServer{
		app:         data.NewReserveData(nil, nil, nil, nil, testStorage, nil, setting),
		core:        core.NewReserveCore(nil, nil, setting),
		metric:      nil,
		authEnabled: false,
		r:           gin.Default(),
		blockchain:  testHTTPBlockchain{},
		setting:     setting,
	}
	testServer.register()
	var tests = []testCase{
		{
			msg: "test update config bad request",
			data: map[string]string{
				"btc": "abc",
			},
			assert: httputil.ExpectFailure(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.msg, func(t *testing.T) { testHTTPRequest(t, tc, testServer.r) })
	}
}
