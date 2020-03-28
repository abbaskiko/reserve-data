package gasstation

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_ETHGas(t *testing.T) {
	c := New(&http.Client{})
	gas, err := c.ETHGas()
	require.NoError(t, err)
	require.True(t, gas.Fast > 0)
	require.True(t, gas.Average > 0)
	require.True(t, gas.Fastest > 0)
	require.True(t, gas.SafeLow > 0)
}
