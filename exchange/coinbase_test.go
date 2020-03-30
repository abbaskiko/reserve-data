package exchange

import (
	"testing"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestCoinbaseInterface(t *testing.T) {
	var cb interface{}
	l, err := zap.NewDevelopment()
	require.NoError(t, err)
	sugar := l.Sugar()
	cb = NewCoinbase(sugar, common.ExchangeID(1), nil, nil)

	// assert coinbase as a common exchange
	_, ok := cb.(common.Exchange)
	assert.True(t, ok)
}
