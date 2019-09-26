package fetcher

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewTickerRunner(t *testing.T) {
	runner := NewTickerRunner(time.Millisecond, time.Millisecond, time.Millisecond, time.Millisecond, time.Millisecond, time.Millisecond)
	err := runner.Start()
	require.NoError(t, err)
	for _, testCase := range []struct {
		ch         <-chan time.Time
		nameTicker string
	}{
		{
			ch:         runner.GetOrderbookTicker(),
			nameTicker: "order book",
		}, {
			ch:         runner.GetAuthDataTicker(),
			nameTicker: "auth data",
		}, {
			ch:         runner.GetRateTicker(),
			nameTicker: "rate",
		}, {
			ch:         runner.GetBlockTicker(),
			nameTicker: "block",
		}, {
			ch:         runner.GetGlobalDataTicker(),
			nameTicker: "global",
		}, {
			ch:         runner.GetExchangeHistoryTicker(),
			nameTicker: "exchange history",
		},
	} {
		select {
		case <-testCase.ch:
			t.Logf("get %v", testCase.nameTicker)
		case <-time.After(time.Millisecond * 2):
			t.Fatalf("ticker %v does not work as expected", testCase.nameTicker)
		}
	}
}
