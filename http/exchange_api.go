package http

import (
	"github.com/gin-gonic/gin"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/http/httputil"
)

type failedCancelOrder struct {
	Reason string            `json:"reason"`
	ID     common.ActivityID `json:"id"`
}

// CancelAllOrders cancel all orders
func (s *Server) CancelAllOrders(c *gin.Context) {
	var (
		response []failedCancelOrder
	)
	pendingActivites, err := s.app.GetPendingActivities()
	if err != nil {
		httputil.ResponseFailure(c, httputil.WithError(err))
		return
	}

	for _, activity := range pendingActivites {
		if activity.Action == common.ActionTrade {
			exchange := activity.Params["exchange"].(common.Exchange)
			// Cancel order
			if err := s.core.CancelOrder(activity.ID, exchange); err != nil {
				// save failed order id
				response = append(response, failedCancelOrder{
					Reason: err.Error(),
					ID:     activity.ID,
				})
				continue
			}
		}
	}

	httputil.ResponseSuccess(c, httputil.WithData(response))
}
