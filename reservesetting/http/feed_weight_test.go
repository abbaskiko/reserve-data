package http

import (
	"testing"

	"github.com/KyberNetwork/reserve-data/reservesetting/common"
	"github.com/KyberNetwork/reserve-data/world"
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
		world.CoinbaseETHBTC.String(): 0.1,
		world.GeminiETHBTC.String():   0.2,
	}
	err = checkFeedWeight(&testSetRateBTCFeed, &testFeedWeightIsSupported)
	assert.NoError(t, err)

	// Test feed weight USD Feed
	testSetRateUSDFeed := common.USDFeed
	err = checkFeedWeight(&testSetRateUSDFeed, &testFeedWeightIsNotSupported)
	require.NotNil(t, err)

	testFeedWeightUSDIsSupported := common.FeedWeight{
		world.CoinbaseETHUSD.String():  0.2,
		world.GeminiETHUSD.String():    0.1,
		world.CoinbaseETHUSDC.String(): 0.3,
		world.BinanceETHUSDC.String():  1.1,
		world.CoinbaseETHDAI.String():  1.2,
		world.HitBTCETHDAI.String():    1.3,
		world.BitFinexETHUSDT.String(): 1.4,
		world.BinanceETHPAX.String():   1.5,
		world.BinanceETHTUSD.String():  0.6,
		world.BinanceETHUSDT.String():  0.7,
	}
	err = checkFeedWeight(&testSetRateUSDFeed, &testFeedWeightUSDIsSupported)
	assert.NoError(t, err)
}
