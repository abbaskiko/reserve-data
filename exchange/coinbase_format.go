package exchange

import (
	"encoding/json"
	"fmt"
)

// CoinbasePrice ...
type CoinbasePrice struct {
	Price    string
	Size     string
	NumOrder int
}

func (n *CoinbasePrice) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{&n.Price, &n.Size, &n.NumOrder}
	wantLen := len(tmp)
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), wantLen; g != e {
		return fmt.Errorf("wrong number of fields in CoinbasePrice: %d != %d", g, e)
	}
	return nil
}

// CoinbaseResp ...
type CoinbaseResp struct {
	Sequence int64           `json:"sequence"`
	Code     int             `json:"code"`
	Msg      string          `json:"msg"`
	Bids     []CoinbasePrice `json:"bids"`
	Asks     []CoinbasePrice `json:"asks"`
}
