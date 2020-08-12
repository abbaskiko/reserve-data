package testutil

import (
	"testing"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/data"
	"github.com/KyberNetwork/reserve-data/data/fetcher"
)

// GlobalStorageTestSuite is a test helper for data.GlobalStorage and fetcher.GlobalStorage.
type GlobalStorageTestSuite struct {
	t   *testing.T
	dgs data.GlobalStorage
	fgs fetcher.GlobalStorage
}

// NewGlobalStorageTestSuite creates new a new test suite with given implementations of data.GlobalStorage and
// fetcher.GlobalStorage.
func NewGlobalStorageTestSuite(t *testing.T, dgs data.GlobalStorage, fgs fetcher.GlobalStorage) *GlobalStorageTestSuite {
	t.Helper()
	return &GlobalStorageTestSuite{
		t:   t,
		dgs: dgs,
		fgs: fgs,
	}
}

// Run executes test suite.
func (ts *GlobalStorageTestSuite) Run() {
	ts.t.Helper()

	gdaxTradeIDs := []uint64{1, 2, 3}
	for _, tradeID := range gdaxTradeIDs {
		err := ts.fgs.StoreGoldInfo(common.GoldData{
			Timestamp: common.NowInMillis(),
			GDAX: common.GDAXGoldData{
				Valid:   true,
				TradeID: tradeID,
			},
		})
		if err != nil {
			ts.t.Fatal(err)
		}
	}

	ts.t.Log("getting gold info of future timepoint, expected to get latest version")
	var futureTimepoint uint64 = common.NowInMillis() + 100
	version, err := ts.dgs.CurrentGoldInfoVersion(futureTimepoint)
	if err != nil {
		ts.t.Fatal(err)
	}

	goldInfo, err := ts.dgs.GetGoldInfo(version)
	if err != nil {
		ts.t.Fatal(err)
	}
	ts.t.Logf("latest gold info version: %d", version)

	lastGdaxTradeID := gdaxTradeIDs[(len(gdaxTradeIDs) - 1)]
	if goldInfo.GDAX.TradeID != lastGdaxTradeID {
		ts.t.Errorf("getting wrong dgx status, expected: %d, got: %d", lastGdaxTradeID, goldInfo.GDAX.TradeID)
	}
}
