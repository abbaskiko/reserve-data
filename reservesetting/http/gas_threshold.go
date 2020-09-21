package http

import (
	"encoding/json"

	"github.com/gin-gonic/gin"

	gaspricedataclient "github.com/KyberNetwork/reserve-data/common/gaspricedata-client"
	"github.com/KyberNetwork/reserve-data/http/httputil"
	"github.com/KyberNetwork/reserve-data/reservesetting/common"
)

type gasStatusResult struct {
	GasPrice gaspricedataclient.GasResult `json:"gas_price"`
	GasThresholdSetting
}

// GasThresholdSetting ...
type GasThresholdSetting struct {
	High float64 `json:"high"`
	Low  float64 `json:"low"`
}

const (
	gasThresholdKey = "gas_threshold" // key in general setting table
)

func (s *Server) getGasStatus(c *gin.Context) {
	gasPrice, err := s.gasClient.GetGas()
	if err != nil {
		s.l.Errorw("query gas price failed", "err", err)
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

	result.GasPrice = gasPrice
	httputil.ResponseSuccess(c, httputil.WithData(result))
}

func (s *Server) setGasThreshold(c *gin.Context) {
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
