package exchange

import (
	"fmt"
	"strings"

	"github.com/KyberNetwork/reserve-data/common"
	commonv3 "github.com/KyberNetwork/reserve-data/reservesetting/common"
)

// HuobiLive is LiveExchange for huobi
type HuobiLive struct {
	interf HuobiInterface
}

// NewHuobiLive return new HuobiLive instance
func NewHuobiLive(interf HuobiInterface) *HuobiLive {
	return &HuobiLive{
		interf: interf,
	}
}

// GetLiveExchangeInfos querry the Exchange Endpoint for exchange precision and limit of a list of tokenPairIDs
// It return error if occurs.
func (hl *HuobiLive) GetLiveExchangeInfos(pairs []commonv3.TradingPairSymbols) (common.ExchangeInfo, error) {
	result := make(common.ExchangeInfo)
	exchangeInfo, err := hl.interf.GetExchangeInfo()
	if err != nil {
		return result, err
	}
	for _, pair := range pairs {
		exchangePrecisionLimit, ok := hl.getPrecisionLimitFromSymbols(pair, exchangeInfo)
		if !ok {
			return result, fmt.Errorf("huobi Exchange Info reply doesn't contain token pair %d, base: %s, quote: %s", pair.ID, pair.BaseSymbol, pair.QuoteSymbol)
		}
		result[pair.ID] = exchangePrecisionLimit
	}
	return result, nil
}

// getPrecisionLimitFromSymbols find the pairID amongs symbols from exchanges,
// return ExchangePrecisionLimit of that pair and true if the pairID exist amongs symbols, false if otherwise
func (hl *HuobiLive) getPrecisionLimitFromSymbols(pair commonv3.TradingPairSymbols, symbols HuobiExchangeInfo) (common.ExchangePrecisionLimit, bool) {
	var result common.ExchangePrecisionLimit
	pairName := strings.ToUpper(fmt.Sprintf("%s%s", pair.BaseSymbol, pair.QuoteSymbol))
	for _, symbol := range symbols.Data {
		symbolName := strings.ToUpper(symbol.Base + symbol.Quote)
		if symbolName == pairName {
			result.Precision.Amount = symbol.AmountPrecision
			result.Precision.Price = symbol.PricePrecision
			result.MinNotional = 0.02
			return result, true
		}
	}
	return result, false
}
