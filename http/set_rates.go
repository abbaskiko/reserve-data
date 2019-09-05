package http

import (
	"fmt"
	"math/big"

	"github.com/gin-gonic/gin"

	"github.com/KyberNetwork/reserve-data/http/httputil"
	v3common "github.com/KyberNetwork/reserve-data/reservesetting/common"
)

// RateRequest is request for a rate
type RateRequest struct {
	AssetID uint64 `json:"asset_id"`
	Buy     string `json:"buy"`
	Sell    string `json:"sell"`
	Mid     string `json:"mid"`
	Msg     string `json:"msg"`
}

// SetRateEntry is input for set rate request
type SetRateEntry struct {
	Block uint64        `json:"block"`
	Rates []RateRequest `json:"rates"`
}

// SetRate is for setting token rate
func (s *Server) SetRate(c *gin.Context) {
	var (
		input     SetRateEntry
		assets    []v3common.Asset
		bigBuys   = []*big.Int{}
		bigSells  = []*big.Int{}
		bigAfpMid = []*big.Int{}
		msgs      []string
	)
	if err := c.ShouldBindJSON(&input); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	for _, rates := range input.Rates {
		asset, err := s.settingStorage.GetAsset(rates.AssetID)
		if err != nil {
			httputil.ResponseFailure(c, httputil.WithError(fmt.Errorf("failed to get asset from storage: %s", err.Error())))
			return
		}
		assets = append(assets, asset)
	}
	for _, rate := range input.Rates {
		rbuy, ok := big.NewInt(0).SetString(rate.Buy, 10)
		if !ok {
			httputil.ResponseFailure(c, httputil.WithError(fmt.Errorf("cannot parse rate number buy: %s", rate.Buy)))
			return
		}
		bigBuys = append(bigBuys, rbuy)
		rSell, ok := big.NewInt(0).SetString(rate.Sell, 10)
		if !ok {
			httputil.ResponseFailure(c, httputil.WithError(fmt.Errorf("cannot parse rate number sell: %s", rate.Sell)))
			return
		}
		bigSells = append(bigSells, rSell)
		rMid, ok := big.NewInt(0).SetString(rate.Mid, 10)
		if !ok {
			httputil.ResponseFailure(c, httputil.WithError(fmt.Errorf("cannot parse rate number mid: %s", rate.Mid)))
			return
		}
		bigAfpMid = append(bigAfpMid, rMid)
		msgs = append(msgs, rate.Msg)
	}
	id, err := s.core.SetRates(assets, bigBuys, bigSells, big.NewInt(int64(input.Block)), bigAfpMid, msgs)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(fmt.Errorf("failed to set rates: %s", err.Error())))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithField("id", id))
}
