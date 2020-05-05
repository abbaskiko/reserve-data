package binance

import (
	"encoding/json"

	"github.com/KyberNetwork/reserve-data/feed-provider/common"
)

type Price struct {
	Price float64 `json:",string"`
	Size  float64 `json:",string"`
}

func (p *Price) UnmarshalJSON(buf []byte) error {
	var price, size json.Number
	tmp := []interface{}{&price, &size}
	err := json.Unmarshal(buf, &tmp)
	if err != nil {
		return err
	}
	p.Price, err = price.Float64()
	if err != nil {
		return err
	}
	p.Size, err = size.Float64()
	if err != nil {
		return err
	}
	return nil
}

func (p *Price) toCommonPrice() common.Price {
	return common.Price{
		Price: p.Price,
		Size:  p.Size,
	}
}

type Resp struct {
	Sequence int64   `json:"sequence"`
	Bids     []Price `json:"bids"`
	Asks     []Price `json:"asks"`
}

func (r *Resp) toOrderbooks() common.Orderbooks {
	var asks, bids []common.Price
	for _, ask := range r.Asks {
		asks = append(asks, ask.toCommonPrice())
	}
	for _, bid := range r.Bids {
		bids = append(bids, bid.toCommonPrice())
	}
	return common.Orderbooks{
		Asks: asks,
		Bids: bids,
	}
}

func (r *Resp) toFeed(amount float64) common.Feed {
	orderbook := r.toOrderbooks()
	return orderbook.ToFeed(amount)
}
