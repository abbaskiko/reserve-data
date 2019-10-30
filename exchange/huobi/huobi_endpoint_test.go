package huobi

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDepositAddress(t *testing.T) {
	t.Skip()     // skip as external test
	key := ""    // enter only once for test
	secret := "" // enter only once for test
	signer := NewSigner(key, secret)
	interf := NewRealInterface()
	ep := NewHuobiEndpoint(*signer, interf)

	depositAddress, err := ep.GetDepositAddress("knc")
	assert.NoError(t, err)

	log.Printf("deposit address: %+v", depositAddress.Data)
}
