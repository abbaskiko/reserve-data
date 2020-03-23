package http

import (
	"testing"

	"github.com/KyberNetwork/reserve-data/common/feed"
	"github.com/KyberNetwork/reserve-data/reservesetting/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFeedWeight(t *testing.T) {

	var (
		testSetRateBTCFeed           = common.BTCFeed
		testFeedWeightIsNotSupported = common.FeedWeight{
			"NotSupportedFeed": 0.341432,
		}
	)
	// Test feed weight is not supported
	err := checkFeedWeight(&testSetRateBTCFeed, &testFeedWeightIsNotSupported)
	require.NotNil(t, err)

	// Test feed weight is supported
	testFeedWeightIsSupported := common.FeedWeight{
		feed.CoinbaseETHBTC.String(): 0.1,
		feed.GeminiETHBTC.String():   0.2,
	}
	err = checkFeedWeight(&testSetRateBTCFeed, &testFeedWeightIsSupported)
	assert.NoError(t, err)

	// Test feed weight USD Feed
	testSetRateUSDFeed := common.USDFeed
	err = checkFeedWeight(&testSetRateUSDFeed, &testFeedWeightIsNotSupported)
	require.NotNil(t, err)

	testFeedWeightUSDIsSupported := common.FeedWeight{
		feed.CoinbaseETHUSD.String():  0.2,
		feed.GeminiETHUSD.String():    0.1,
		feed.CoinbaseETHUSDC.String(): 0.3,
		feed.BinanceETHUSDC.String():  1.1,
		feed.CoinbaseETHDAI.String():  1.2,
		feed.HitBTCETHDAI.String():    1.3,
		feed.BitFinexETHUSDT.String(): 1.4,
		feed.BinanceETHPAX.String():   1.5,
		feed.BinanceETHTUSD.String():  0.6,
		feed.BinanceETHUSDT.String():  0.7,
	}
	err = checkFeedWeight(&testSetRateUSDFeed, &testFeedWeightUSDIsSupported)
	assert.NoError(t, err)
}
