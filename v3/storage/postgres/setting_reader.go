package postgres

import (
	"database/sql"
	"log"

	ethereum "github.com/ethereum/go-ethereum/common"

	"github.com/KyberNetwork/reserve-data/v3/common"
)

// TODO: rewritten this function to filter the set rate strategy in SQL query
func (s *Storage) GetTransferableAssets() ([]common.Asset, error) {
	var results []common.Asset
	allAssets, err := s.GetAssets()
	if err != nil {
		return nil, err
	}

	for _, asset := range allAssets {
		if asset.SetRate != common.SetRateNotSet {
			results = append(results, asset)
		}
	}

	return results, nil
}

type tradingPairSymbolsDB struct {
	BaseSymbol  string `db:"base_symbol"`
	QuoteSymbol string `db:"quote_symbol"`
}

func (s *Storage) GetTradingPairSymbols(id uint64) ([]common.TradingPairSymbols, error) {
	var (
		tradingPairs []tradingPairSymbolsDB
		result       []common.TradingPairSymbols
	)
	if err := s.stmts.getTradingPairSymbols.Select(&tradingPairs, id); err != nil {
		return nil, err
	}
	for _, pair := range tradingPairs {
		result = append(result, common.TradingPairSymbols{
			BaseSymbol:  pair.BaseSymbol,
			QuoteSymbol: pair.QuoteSymbol,
		})
	}
	return result, nil
}

// TODO: rewrite this function to filter the exchange if in SQL query.
func (s *Storage) GetDepositAddresses(exchangeID uint64) (map[string]ethereum.Address, error) {
	var results = make(map[string]ethereum.Address)

	allAssets, err := s.GetAssets()
	if err != nil {
		return nil, err
	}

	for _, asset := range allAssets {
		for _, exchange := range asset.Exchanges {
			if exchange.ExchangeID == exchangeID {
				results[exchange.Symbol] = exchange.DepositAddress
			}
		}
	}

	return results, nil
}

func (s *Storage) GetMinNotional(exchangeID, baseID, quoteID uint64) (float64, error) {
	var minNotional float64
	log.Printf("getting min notional for exchange=%d base=%d quote=%d",
		exchangeID, baseID, quoteID)
	if err := s.stmts.getMinNotional.Get(&minNotional,
		exchangeID, baseID, quoteID); err == sql.ErrNoRows {
		return 0, common.ErrNotFound
	} else if err != nil {
		return 0, err
	}
	return minNotional, nil
}

// TODO: rewrite this function with proper SQL statement
func (s *Storage) GetAssetBySymbol(exchangeID uint64, symbol string) (common.Asset, error) {
	allAssets, err := s.GetAssets()
	if err != nil {
		return common.Asset{}, err
	}

	for _, asset := range allAssets {
		for _, exchange := range asset.Exchanges {
			if exchange.ExchangeID == exchangeID && exchange.Symbol == symbol {
				return asset, nil
			}
		}
	}

	return common.Asset{}, common.ErrNotFound
}
