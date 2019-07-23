package common

import (
	"encoding/json"
	"fmt"
	"time"

	ethereum "github.com/ethereum/go-ethereum/common"
)

// Exchange represents a centralized exchange in database.
type Exchange struct {
	ID              uint64  `json:"id"`
	Name            string  `json:"name"`
	TradingFeeMaker float64 `json:"trading_fee_maker"`
	TradingFeeTaker float64 `json:"trading_fee_taker"`
	Disable         bool    `json:"disable"`
}

//go:generate stringer -type=SetRate -linecomment
type SetRate int

const (
	// SetRateNotSet indicates that set rate is not enabled for the asset.
	SetRateNotSet SetRate = iota // not_set
	// ExchangeFeed is used when asset rate is set from fetching historical prices from exchanges.
	ExchangeFeed // exchange_feed
	// GoldFeed is used when asset rate is set from fetching gold prices.
	GoldFeed // gold_feed
	// BTCFeed is used when asset rate is set from fetching Bitcoin prices.
	BTCFeed // btc_feed
)

var validSetRateTypes = map[string]SetRate{
	SetRateNotSet.String(): SetRateNotSet,
	ExchangeFeed.String():  ExchangeFeed,
	GoldFeed.String():      GoldFeed,
	BTCFeed.String():       BTCFeed,
}

// SetRateFromString returns the SetRate value from its string presentation, if exists.
func SetRateFromString(s string) (SetRate, bool) {
	sr, ok := validSetRateTypes[s]
	return sr, ok
}

func (i SetRate) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, i.String())), nil
}
func isString(input []byte) bool {
	return len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"'
}
func (i *SetRate) UnmarshalJSON(input []byte) error {
	if !isString(input) {
		return fmt.Errorf("not is string")
	}
	r, ok := SetRateFromString(string(input[1 : len(input)-1]))
	if !ok {
		return fmt.Errorf("%s is not a valid SetRate", input)
	}
	*i = r
	return nil
}

// TradingPair is a trading in an exchange.
type TradingPair struct {
	ID              uint64  `json:"id"`
	Base            uint64  `json:"base"`
	Quote           uint64  `json:"quote"`
	PricePrecision  uint64  `json:"price_precision"`
	AmountPrecision uint64  `json:"amount_precision"`
	AmountLimitMin  float64 `json:"amount_limit_min"`
	AmountLimitMax  float64 `json:"amount_limit_max"`
	PriceLimitMin   float64 `json:"price_limit_min"`
	PriceLimitMax   float64 `json:"price_limit_max"`
	MinNotional     float64 `json:"min_notional"`
}

type TradingPairSymbols struct {
	TradingPair
	BaseSymbol  string `json:"base_symbol"`
	QuoteSymbol string `json:"quote_symbol"`
}

// AssetExchange is the configuration of an asset for a specific exchange.
type AssetExchange struct {
	ID                uint64           `json:"id"`
	ExchangeID        uint64           `json:"exchange_id"`
	Symbol            string           `json:"symbol"`
	DepositAddress    ethereum.Address `json:"deposit_address"`
	MinDeposit        float64          `json:"min_deposit"`
	WithdrawFee       float64          `json:"withdraw_fee"`
	TargetRecommended float64          `json:"target_recommended"`
	TargetRatio       float64          `json:"target_ratio"`
	TradingPairs      []TradingPair    `json:"trading_pairs" binding:"dive"`
}

// AssetTarget is the target setting of an asset.
type AssetTarget struct {
	Total              float64 `json:"total"`
	Reserve            float64 `json:"reserve"`
	RebalanceThreshold float64 `json:"rebalance_threshold"`
	TransferThreshold  float64 `json:"transfer_threshold"`
}

// PWIEquation is a PWI equation. An asset will have 2 PWI equation: ask and bid.
type PWIEquation struct {
	A                   float64 `json:"a"`
	B                   float64 `json:"b"`
	C                   float64 `json:"c"`
	MinMinSpread        float64 `json:"min_min_spread"`
	PriceMultiplyFactor float64 `json:"price_multiply_factor"`
}

// AssetPWI is the PWI configuration of an asset.
type AssetPWI struct {
	Ask PWIEquation `json:"ask"`
	Bid PWIEquation `json:"bid"`
}

type RebalanceQuadratic struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
	C float64 `json:"c"`
}

// Asset represents an asset in centralized exchange, eg: ETH, KNC, Bitcoin...
type Asset struct {
	ID                 uint64              `json:"id"`
	Symbol             string              `json:"symbol"`
	Name               string              `json:"name"`
	Address            ethereum.Address    `json:"address"`
	OldAddresses       []ethereum.Address  `json:"old_addresses"`
	Decimals           uint64              `json:"decimals"`
	Transferable       bool                `json:"transferable"`
	SetRate            SetRate             `json:"set_rate"`
	Rebalance          bool                `json:"rebalance"`
	IsQuote            bool                `json:"is_quote"`
	PWI                *AssetPWI           `json:"pwi"`
	RebalanceQuadratic *RebalanceQuadratic `json:"rebalance_quadratic"`
	Exchanges          []AssetExchange     `json:"exchanges" binding:"dive"`
	Target             *AssetTarget        `json:"target"`
	Created            time.Time           `json:"created"`
	Updated            time.Time           `json:"updated"`
}

// TODO: write custom marshal json for created/updated fields

// CreateAsset hold state of being create Asset and waiting for confirm to be Asset.
type CreateAsset struct {
	ID      uint64          `json:"id"`
	Created time.Time       `json:"created"`
	Data    json.RawMessage `json:"data"`
}

// CreateAssetExchange is the configuration of an asset for a specific exchange.
type CreateAssetExchange struct {
	AssetID           uint64           `json:"asset_id"`
	ExchangeID        uint64           `json:"exchange_id"`
	Symbol            string           `json:"symbol"`
	DepositAddress    ethereum.Address `json:"deposit_address"`
	MinDeposit        float64          `json:"min_deposit"`
	WithdrawFee       float64          `json:"withdraw_fee"`
	TargetRecommended float64          `json:"target_recommended"`
	TargetRatio       float64          `json:"target_ratio"`
}

// CreateAssetEntry represents an asset in centralized exchange, eg: ETH, KNC, Bitcoin...
type CreateAssetEntry struct {
	Symbol             string              `json:"symbol" binding:"required"`
	Name               string              `json:"name" binding:"required"`
	Address            ethereum.Address    `json:"address"`
	OldAddresses       []ethereum.Address  `json:"old_addresses"`
	Decimals           uint64              `json:"decimals"`
	Transferable       bool                `json:"transferable"`
	SetRate            SetRate             `json:"set_rate"`
	Rebalance          bool                `json:"rebalance"`
	IsQuote            bool                `json:"is_quote"`
	PWI                *AssetPWI           `json:"pwi"`
	RebalanceQuadratic *RebalanceQuadratic `json:"rebalance_quadratic"`
	Exchanges          []AssetExchange     `json:"exchanges" binding:"dive"`
	Target             *AssetTarget        `json:"target"`
}

// CreateCreateAsset present for a CreateAsset(pending) request
type CreateCreateAsset struct {
	AssetInputs []CreateAssetEntry `json:"assets" binding:"required,dive"`
}

// UpdateAsset hold state of being update Asset and waiting for confirm to apply.
type UpdateAsset struct {
	ID      uint64          `json:"id"`
	Created time.Time       `json:"created"`
	Data    json.RawMessage `json:"data"`
}

// CreateUpdateAsset present for an UpdateAsset request
type CreateUpdateAsset struct {
	AssetID      uint64            `json:"asset_id"`
	Symbol       *string           `json:"symbol"`
	Name         *string           `json:"name"`
	Address      *ethereum.Address `json:"address"`
	Decimals     *uint64           `json:"decimals"`
	Transferable *bool             `json:"transferable"`
	SetRate      *SetRate          `json:"set_rate"`
	Rebalance    *bool             `json:"rebalance"`
	IsQuote      *bool             `json:"is_quote"`
}

// UpdateAssetExchange is the options of UpdateAssetExchange method.
type UpdateAssetExchange struct {
	Symbol            *string           `json:"symbol"`
	DepositAddress    *ethereum.Address `json:"deposit_address"`
	MinDeposit        *float64          `json:"min_deposit"`
	WithdrawFee       *float64          `json:"withdraw_fee"`
	TargetRecommended *float64          `json:"target_recommended"`
	TargetRatio       *float64          `json:"target_ratio"`
}

type UpdateExchange struct {
	TradingFeeMaker *float64 `json:"trading_fee_maker"`
	TradingFeeTaker *float64 `json:"trading_fee_taker"`
	Disable         *bool    `json:"disable"`
}
