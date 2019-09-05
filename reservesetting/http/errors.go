package http

import (
	"errors"

	"github.com/KyberNetwork/reserve-data/reservesetting/common"
)

var errorMapping = map[error]error{
	common.ErrAssetExchangeMissing:      errors.New("rebalance is enabled, will require exchanges define for the asset"),
	common.ErrAssetTargetMissing:        errors.New("rebalance is enabled, but target configuration is not set"),
	common.ErrRebalanceQuadraticMissing: errors.New("rebalance is enabled, but rebalance quadratic is not set"),
}

func makeFriendlyMessage(err error) error {
	// we try to make a friendly error message from app error, eg 'missing asset exchange configuration'
	// can confuse to human for what's missing.
	if newErr, ok := errorMapping[err]; ok {
		return newErr
	}
	return err
}
