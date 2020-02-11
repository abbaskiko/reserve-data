package exchange

import (
	"errors"
)

const (
	tradeTypeBuy       = "buy"
	tradeTypeSell      = "sell"
	exchangeStatusDone = "done"
)

var ErrNotSupport = errors.New("not support")
