package coinbase

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-data/exchange"
)

// Endpoint endpoint object
type Endpoint struct {
	interf Interface
	l      *zap.SugaredLogger
	client *http.Client
}

// GetResponse call to binance endpoint and get response
func (ep *Endpoint) GetResponse(method string, url string, params map[string]string) ([]byte, error) {
	var (
		err      error
		respBody []byte
	)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	ep.l.Infof("request to coinbase: %s", req.URL)
	resp, err := ep.client.Do(req)
	if err != nil {
		return respBody, err
	}
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return respBody, err
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return respBody, fmt.Errorf("coinbase return not OK code %d", resp.StatusCode)
	}
	return respBody, err
}

func (ep *Endpoint) GetOnePairOrderBook(baseID, quoteID string) (exchange.CoinbaseResp, error) {
	respBody, err := ep.GetResponse(
		"GET",
		ep.interf.PublicEndpoint()+fmt.Sprintf("/products/%s-%s/book?level=2", baseID, quoteID), nil)

	respData := exchange.CoinbaseResp{}
	if err != nil {
		return respData, err
	}
	if err = json.Unmarshal(respBody, &respData); err != nil {
		return respData, err
	}
	if respData.Code != 0 {
		return respData, fmt.Errorf("getting orderbook from coinbase failed: %s", respData.Msg)
	}
	return respData, nil
}

//NewHuobiEndpoint return new endpoint instance
func NewCoinbaseEndpoint(interf Interface, client *http.Client) *Endpoint {
	return &Endpoint{
		interf: interf,
		l:      zap.S(),
		client: client,
	}
}
