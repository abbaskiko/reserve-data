package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	v1common "github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/common/testutil"
	"github.com/KyberNetwork/reserve-data/reservesetting/common"
	"github.com/KyberNetwork/reserve-data/reservesetting/storage/postgres"
	"github.com/KyberNetwork/reserve-data/world"
)

func TestCheckFeedConfiguration(t *testing.T) {
	db, tearDown := testutil.MustNewDevelopmentDB()
	defer func() {
		assert.NoError(t, tearDown())
	}()
	s, err := postgres.NewStorage(db)
	require.NoError(t, err)
	supportedExchanges := make(map[v1common.ExchangeID]v1common.LiveExchange)
	server := NewServer(s, "", supportedExchanges, "", "")

	feedConfigurationEntryNotSupportedTest := common.SetFeedConfigurationEntry{
		Name: "Not supported feed",
	}
	err = server.checkSetFeedConfigurationParams(feedConfigurationEntryNotSupportedTest)
	require.Error(t, err)

	// test for supported feed names
	supportedFeed := []world.Feed{
		world.GoldData,
		world.OneForgeXAUETH,
		world.OneForgeXAUUSD,
		world.GDAXETHUSD,
		world.KrakenETHUSD,
		world.GeminiETHUSD,
		world.CoinbaseETHBTC,
		world.GeminiETHBTC,
		world.CoinbaseETHUSDC,
		world.BinanceETHUSDC,
		world.CoinbaseETHUSD,
		world.CoinbaseETHDAI,
		world.HitBTCETHDAI,
		world.BitFinexETHUSDT,
		world.BinanceETHUSDT,
		world.BinanceETHPAX,
		world.BinanceETHTUSD,
	}

	for _, feedName := range supportedFeed {
		feedNameSupportedEntry := common.SetFeedConfigurationEntry{
			Name: feedName.String(),
		}
		err := server.checkSetFeedConfigurationParams(feedNameSupportedEntry)
		assert.NoError(t, err)
	}
}
