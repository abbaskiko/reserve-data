package huobi

import (
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	huobiEndpoint = "https://api.huobi.pro"
)

func TestGetDepositAddress(t *testing.T) {
	t.Skip()     // skip as external test
	key := ""    // enter only once for test
	secret := "" // enter only once for test
	signer := NewSigner(key, secret)
	interf := NewRealInterface(huobiEndpoint)
	ep := NewHuobiEndpoint(signer, interf, &http.Client{Timeout: time.Second * 30})

	depositAddress, err := ep.GetDepositAddress("ETH")
	assert.NoError(t, err)

	log.Printf("deposit address: %+v", depositAddress.Data)
}
