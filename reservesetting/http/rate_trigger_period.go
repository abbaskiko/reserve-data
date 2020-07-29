package http

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/KyberNetwork/reserve-data/http/httputil"
	"github.com/KyberNetwork/reserve-data/reservesetting/common"
)

const (
	rateTriggerPeriodKey = "rate_trigger_period"
)

type rateTriggerPeriod map[string]float64

func (s *Server) setRateTriggerPeriod(c *gin.Context) {
	var input struct {
		Value float64 `json:"value" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	if input.Value <= 0 {
		httputil.ResponseFailure(c, httputil.WithError(fmt.Errorf("value must greater than zero, value=%f", input.Value)))
		return
	}
	gdata := common.GeneralData{
		Key:   rateTriggerPeriodKey,
		Value: strconv.FormatFloat(input.Value, 'f', -1, 64),
	}
	id, err := s.storage.SetGeneralData(gdata)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithData(id))
}

func (s *Server) getRateTriggerPeriod(c *gin.Context) {
	gdata, err := s.storage.GetGeneralData(rateTriggerPeriodKey)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	value, err := strconv.ParseFloat(gdata.Value, 64)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithData(rateTriggerPeriod{
		rateTriggerPeriodKey: value,
	}))
}
