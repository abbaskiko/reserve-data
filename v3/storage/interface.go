package storage

import (
	ethereum "github.com/ethereum/go-ethereum/common"

	v3 "github.com/KyberNetwork/reserve-data/v3/common"
)

// Interface is the common persistent storage interface of V3 APIs.
type Interface interface {
	SettingReader

	GetExchanges() ([]v3.Exchange, error)
	UpdateExchange(id uint64, opts UpdateExchangeOpts) error

	CreateAssetExchange(exchangeID, assetID uint64, symbol string, depositAddress ethereum.Address,
		minDeposit, withdrawFee, targetRecommended, targetRatio float64) (uint64, error)
	UpdateAssetExchange(id uint64, opts UpdateAssetExchangeOpts) error

	CreateCreateAssetExchange(v3.CreateCreateAssetExchange) (uint64, error)
	GetCreateAssetExchanges() ([]v3.CreateAssetExchange, error)
	GetCreateAssetExchange(id uint64) (v3.CreateAssetExchange, error)
	RejectCreateAssetExchange(id uint64) error
	ConfirmCreateAssetExchange(id uint64) error

	CreateUpdateAssetExchange(v3.CreateUpdateAssetExchange) (uint64, error)
	GetUpdateAssetExchanges() ([]v3.UpdateAssetExchange, error)
	GetUpdateAssetExchange(id uint64) (v3.UpdateAssetExchange, error)
	RejectUpdateAssetExchange(id uint64) error
	ConfirmUpdateAssetExchange(id uint64) error

	CreateAsset(
		symbol, name string,
		address ethereum.Address,
		decimals uint64,
		transferable bool,
		setRate v3.SetRate,
		rebalance bool,
		isQuote bool,
		pwi *v3.AssetPWI,
		rb *v3.RebalanceQuadratic,
		exchanges []v3.AssetExchange,
		target *v3.AssetTarget,
	) (uint64, error)
	UpdateAsset(id uint64, opts UpdateAssetOpts) error
	// ChangeAssetAddress make the current address address of asset old address and set new address as current.
	ChangeAssetAddress(id uint64, address ethereum.Address) error
	UpdateDepositAddress(assetID, exchangeID uint64, address ethereum.Address) error

	UpdateTradingPair(id uint64, opts UpdateTradingPairOpts) error

	// TODO method for batch update PWI
	// TODO method for batch update rebalance quadratic
	// TODO method for batch update exchange configuration
	// TODO meethod for batch update target
	// TODO method for update address
	CreateCreateAsset(v3.CreateCreateAsset) (uint64, error)
	GetCreateAssets() ([]v3.CreateAsset, error)
	GetCreateAsset(uint64) (v3.CreateAsset, error)
	RejectCreateAsset(id uint64) error
	ConfirmCreateAsset(id uint64) error

	CreateUpdateAsset(asset v3.CreateUpdateAsset) (uint64, error)
	GetUpdateAssets() ([]v3.UpdateAsset, error)
	GetUpdateAsset(uint64) (v3.UpdateAsset, error)
	RejectUpdateAsset(id uint64) error
	ConfirmUpdateAsset(id uint64) error

	CreateUpdateExchange(exchange v3.CreateUpdateExchange) (uint64, error)
	GetUpdateExchanges() ([]v3.UpdateExchange, error)
	GetUpdateExchange(id uint64) (v3.UpdateExchange, error)
	RejectUpdateExchange(id uint64) error
	ConfirmUpdateExchange(id uint64) error

	CreateCreateTradingPair(v3.CreateCreateTradingPair) (uint64, error)
	GetCreateTradingPairs() ([]v3.CreateTradingPair, error)
	GetCreateTradingPair(id uint64) (v3.CreateTradingPair, error)
	RejectCreateTradingPair(id uint64) error
	ConfirmCreateTradingPair(id uint64) error

	CreateUpdateTradingPair(v3.CreateUpdateTradingPair) (uint64, error)
	GetUpdateTradingPairs() ([]v3.UpdateTradingPair, error)
	GetUpdateTradingPair(id uint64) (v3.UpdateTradingPair, error)
	RejectUpdateTradingPair(id uint64) error
	ConfirmUpdateTradingPair(id uint64) error
}

// SettingReader is the common interface for reading exchanges, assets configuration.
type SettingReader interface {
	GetAsset(id uint64) (v3.Asset, error)
	GetExchange(id uint64) (v3.Exchange, error)
	// TODO: add GetTradingPair method that accept trading_pair_id
	GetTradingPairs(exchangeID uint64) ([]v3.TradingPairSymbols, error)
	// TODO: check usages of this method to see if it should be replaced with GetDepositAddress(exchangeID, tokenID)
	GetDepositAddresses(exchangeID uint64) (map[string]ethereum.Address, error)
	GetAssets() ([]v3.Asset, error)
	// GetTransferableAssets returns all assets that the set rate strategy is not not_set.
	GetTransferableAssets() ([]v3.Asset, error)
	GetMinNotional(exchangeID, baseID, quoteID uint64) (float64, error)
}

// these type match user type define in common package so we just need to make an alias here
// in case they did not later, we need to redefine the structure here, or review this again.
type UpdateAssetExchangeOpts = v3.UpdateAssetExchangeEntry

type UpdateExchangeOpts struct {
	TradingFeeMaker *float64 `json:"trading_fee_maker"`
	TradingFeeTaker *float64 `json:"trading_fee_taker"`
	Disable         *bool    `json:"disable"`
}

type UpdateAssetOpts struct {
	Symbol             *string                `json:"symbol"`
	Name               *string                `json:"name"`
	Address            *ethereum.Address      `json:"address"`
	Decimals           *uint64                `json:"decimals"`
	Transferable       *bool                  `json:"transferable"`
	SetRate            *v3.SetRate            `json:"set_rate"`
	Rebalance          *bool                  `json:"rebalance"`
	IsQuote            *bool                  `json:"is_quote"`
	PWI                *v3.AssetPWI           `json:"pwi"`
	RebalanceQuadratic *v3.RebalanceQuadratic `json:"rebalance_quadratic"`
	Target             *v3.AssetTarget        `json:"target"`
}

// UpdateTradingPairOpts ...
type UpdateTradingPairOpts = v3.UpdateTradingPairEntry
