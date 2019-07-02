package common

import "errors"

var (
	// ErrNotFound is the error to return when no record is found in database.
	ErrNotFound = errors.New("not found")
	// ErrAddressMissing is returned when the operation require address, but not provided or zero.
	ErrAddressMissing = errors.New("address is zero or missing")
	// ErrSymbolExists is returned when creating new asset with duplicated symbol.
	ErrSymbolExists = errors.New("symbol already exists")
	// ErrAddressExists is returned when the address to create is already exists.
	ErrAddressExists = errors.New("address already exists")
	// ErrExchangeFeeMissing is the error to return when user try to enable exchange, but fees are not set.
	ErrExchangeFeeMissing = errors.New("missing exchange fee configuration")
	// ErrPWIMissing is returned when PWI configuration is missing when set rate strategy is defined
	ErrPWIMissing = errors.New("missing PWI configuration")
	// ErrRebalanceQuadraticMissing is returned when rebalance quadratic configuration is missing when
	// rebalance is set to true.
	ErrRebalanceQuadraticMissing = errors.New("missing rebalance quadratic configuration")
	// ErrAssetExchangeMissing is returned when asset exchange configuration is missing for asset with
	// rebalance set to true.
	ErrAssetExchangeMissing = errors.New("missing asset exchange configuration")
	// ErrAssetTargetMissing is returned then asset target configuration is missing for asset with
	// rebalance set to true.
	ErrAssetTargetMissing = errors.New("missing asset target configuration")
	// ErrBadTradingPairConfiguration is returned when bad trading pair configuration is given.
	ErrBadTradingPairConfiguration = errors.New("bad trading pair configuration")
)
