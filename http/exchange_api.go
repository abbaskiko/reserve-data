package http

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/KyberNetwork/reserve-data/common"
	"github.com/KyberNetwork/reserve-data/http/httputil"
)

type failedCancelOrder struct {
	Reason string `json:"reason"`
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
	cancelOrders := make(map[common.Exchange][]common.RequestOrder)
	for _, activity := range pendingActivites {
		if activity.Action == common.ActionTrade {
			exchangeID := activity.Params.Exchange
			// Cancel order
			exchange, ok := common.SupportedExchanges[exchangeID]
			if !ok {
				httputil.ResponseFailure(c, httputil.WithError(fmt.Errorf("exchange %s does not exist", exchange.ID().String())))
				return
			}
			cancelOrders[exchange] = append(cancelOrders[exchange],
				common.RequestOrder{
					ID:     activity.EID,
					Symbol: activity.Params.Base + activity.Params.Quote,
				})
		}
	}
	for exchange, orders := range cancelOrders {
		result := s.core.CancelOrders(orders, exchange)
		for id, res := range result {
			if !res.Success {
				// save failed order id
				response = append(response, failedCancelOrder{
					Reason: fmt.Sprintf("exchange: %s, order: %s, err: %s", exchange.ID().String(), id, err),
				})
			}
		}
	}

	httputil.ResponseSuccess(c, httputil.WithData(response))
}
