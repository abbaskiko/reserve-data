package http

import (
	"encoding/json"

	"github.com/gin-gonic/gin"

	"github.com/KyberNetwork/reserve-data/http/httputil"
	"github.com/KyberNetwork/reserve-data/reservesetting/common"
)

type GsGas struct {
	Fast     float64 `json:"fast"`
	Standard float64 `json:"standard"`
	Slow     float64 `json:"slow"`
}
type gasStatusResult struct {
	GasStation GsGas `json:"eth_gas_station"`
	GasThresholdSetting
}
type GasThresholdSetting struct {
	High float64 `json:"high"`
	Low  float64 `json:"low"`
}

const (
	gasThresholdKey = "gas_threshold" // key in general setting table
)

func (s *Server) GetGasStatus(c *gin.Context) {
	gasPrice, err := s.gasStation.ETHGas()
	if err != nil {
		s.l.Errorw("query gasstation failed", "err", err)
		httputil.ResponseFailure(c, httputil.WithReason(err.Error()))
		return
	}
	gasThresholdData, err := s.storage.GetGeneralData(gasThresholdKey)
	if err != nil {
		s.l.Errorw("failed to get gas threshold", "err", err)
		httputil.ResponseFailure(c, httputil.WithReason(err.Error()))
		return
	}

	result := gasStatusResult{}
	err = json.Unmarshal([]byte(gasThresholdData.Value), &result.GasThresholdSetting)
	if err != nil {
		s.l.Errorw("failed to get gas threshold", "err", err)
		httputil.ResponseFailure(c, httputil.WithReason(err.Error()))
		return
	}

	result.GasStation = GsGas{
		Fast:     gasPrice.Fast / 10.0,
		Standard: gasPrice.Average / 10.0,
		Slow:     gasPrice.SafeLow / 10.0,
	}
	httputil.ResponseSuccess(c, httputil.WithData(result))
}

func (s *Server) SetGasThreshold(c *gin.Context) {
	var v GasThresholdSetting
	if err := c.BindJSON(&v); err != nil {
		httputil.ResponseFailure(c, httputil.WithReason("invalid high-low value"))
		return
	}
	if v.Low >= v.High {
		httputil.ResponseFailure(c, httputil.WithReason("high must > low value"))
		return
	}
	data, err := json.Marshal(v)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithReason(err.Error()))
		return
	}
	if _, err := s.storage.SetGeneralData(common.GeneralData{Key: gasThresholdKey, Value: string(data)}); err != nil {
		httputil.ResponseFailure(c, httputil.WithReason(err.Error()))
		return
	}
	httputil.ResponseSuccess(c)
}
