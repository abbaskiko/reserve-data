package http

import (
	"encoding/json"
	"fmt"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/http/httputil"
	"github.com/gin-gonic/gin"
)

//CheckRebalanceQuadraticRequest check if request data is valid
//rq (requested data) follow format map["tokenID"]{"a": float64, "b": float64, "c": float64}
func (s *Server) CheckRebalanceQuadraticRequest(rq common.RebalanceQuadraticRequest) error {
	for tokenID := range rq {
		if _, err := s.setting.GetInternalTokenByID(tokenID); err != nil {
			return fmt.Errorf("getting token %s got err %s", tokenID, err.Error())
		}
	}
	return nil
}

//SetRebalanceQuadratic set pending rebalance quadratic equation
//input data follow json: {"data":{"KNC": {"a": 0.7, "b": 1.2, "c": 1.3}}}
func (s *Server) SetRebalanceQuadratic(c *gin.Context) {
	postForm, ok := s.Authenticated(c, []string{"value"}, []Permission{ConfigurePermission})
	if !ok {
		return
	}
	value := []byte(postForm.Get("value"))
	if len(value) > maxDataSize {
		httputil.ResponseFailure(c, httputil.WithReason(errDataSizeExceed.Error()))
		return
	}
	var rq common.RebalanceQuadraticRequest
	if err := json.Unmarshal(value, &rq); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	if err := s.CheckRebalanceQuadraticRequest(rq); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	if err := s.metric.StorePendingRebalanceQuadratic(value); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c)
}

//GetPendingRebalanceQuadratic return currently pending config for rebalance quadratic equation
//if there is no pending equation return success false
func (s *Server) GetPendingRebalanceQuadratic(c *gin.Context) {
	_, ok := s.Authenticated(c, []string{}, []Permission{ReadOnlyPermission, ConfigurePermission, ConfirmConfPermission, RebalancePermission})
	if !ok {
		return
	}

	data, err := s.metric.GetPendingRebalanceQuadratic()
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithData(data))
}

//ConfirmRebalanceQuadratic confirm configuration for current pending config for rebalance quadratic equation
func (s *Server) ConfirmRebalanceQuadratic(c *gin.Context) {
	postForm, ok := s.Authenticated(c, []string{}, []Permission{ConfirmConfPermission})
	if !ok {
		return
	}
	value := []byte(postForm.Get("value"))
	if len(value) > maxDataSize {
		httputil.ResponseFailure(c, httputil.WithReason(errDataSizeExceed.Error()))
		return
	}
	err := s.metric.ConfirmRebalanceQuadratic(value)
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c)
}

//RejectRebalanceQuadratic reject pending configuration for rebalance quadratic function
func (s *Server) RejectRebalanceQuadratic(c *gin.Context) {
	_, ok := s.Authenticated(c, []string{}, []Permission{ConfirmConfPermission})
	if !ok {
		return
	}
	if err := s.metric.RemovePendingRebalanceQuadratic(); err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c)
}

//GetRebalanceQuadratic return current confirmed rebalance quadratic equation
func (s *Server) GetRebalanceQuadratic(c *gin.Context) {
	_, ok := s.Authenticated(c, []string{}, []Permission{ReadOnlyPermission, ConfigurePermission, ConfirmConfPermission, RebalancePermission})
	if !ok {
		return
	}

	data, err := s.metric.GetRebalanceQuadratic()
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}
	httputil.ResponseSuccess(c, httputil.WithData(data))
}
