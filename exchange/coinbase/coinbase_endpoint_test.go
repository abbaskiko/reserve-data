package coinbase

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	coinbaseEndpoint = "https://api.pro.coinbase.com"
)

func TestGetPrices(t *testing.T) {
	// t.Skip() // skip as external test
	interf := NewRealInterface(coinbaseEndpoint)
	ep := NewCoinbaseEndpoint(interf, &http.Client{})

	prices, err := ep.GetOnePairOrderBook("eth", "dai")
	assert.NoError(t, err)
	t.Logf("%+v", prices)
}
