package common

import (
	"fmt"
	"sort"
)

type Price struct {
	Price float64 `json:",string"`
	Size  float64 `json:",string"`
}

type Orderbooks struct {
	Asks []Price
	Bids []Price
}

func (o *Orderbooks) ToFeed(amount float64) Feed {
	ask, err := getAfp(o.Asks, amount, false)
	if err != nil {
		return Feed{
			Error: err.Error(),
		}
	}
	bid, err := getAfp(o.Bids, amount, true)
	if err != nil {
		return Feed{
			Error: err.Error(),
		}
	}
	return Feed{
		Ask:   ask,
		Bid:   bid,
		Valid: true,
	}
}

func getAfp(orderbooks []Price, amount float64, isBid bool) (float64, error) {
	var (
		remain                = amount
		totalSize, totalValue float64
	)
	sort.Slice(orderbooks, func(p, q int) bool {
		if isBid {
			return orderbooks[p].Price > orderbooks[q].Price
		}
		return orderbooks[p].Price < orderbooks[q].Price
	})
	// read orderbook from mid until require size is meet
	// and calculate avg_price = sum_product(price, size)/sum(size)
	for _, o := range orderbooks {
		size := o.Size
		if remain < o.Price*size {
			size = remain / o.Price
		}
		value := o.Price * size
		totalValue += value
		totalSize += size
		remain -= value
		if remain <= 0 {
			break
		}
	}
	if remain > 0 {
		return 0, fmt.Errorf("could not collect %f, isBid is %v", amount, isBid)
	}
	if totalSize == 0 {
		return 0, fmt.Errorf("total size of orderbook is 0, isBid is %v", isBid)
	}
	return totalValue / totalSize, nil
}
