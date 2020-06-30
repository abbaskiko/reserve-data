package gasstation

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ETHGas ...
type ETHGas struct {
	Fast          float64            `json:"fast"`
	Fastest       float64            `json:"fastest"`
	SafeLow       float64            `json:"safeLow"`
	Average       float64            `json:"average"`
	BlockTime     float64            `json:"block_time"`
	BlockNum      uint64             `json:"blockNum"`
	Speed         float64            `json:"speed"`
	SafeLowWait   float64            `json:"safeLowWait"`
	AvgWait       float64            `json:"avgWait"`
	FastestWait   float64            `json:"fastestWait"`
	GasPriceRange map[string]float64 `json:"gasPriceRange"`
}

// Client represent for gasStation client
type Client struct {
	client  *http.Client
	baseURL string
	apiKey  string
}

// New create a new Client object
func New(c *http.Client, apiKey string) *Client {
	return &Client{
		client:  c,
		baseURL: "https://ethgasstation.info",
		apiKey:  apiKey,
	}
}

func (c *Client) doRequest(method, path string, response interface{}) error {
	req, err := http.NewRequest(method, c.baseURL+path, nil)
	if err != nil {
		return err
	}
	if c.apiKey != "" {
		q := req.URL.Query()
		q.Set("api-key", c.apiKey)
		req.URL.RawQuery = q.Encode()
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("gasstation return code %d, data %s", resp.StatusCode, string(data))
	}
	err = json.Unmarshal(data, response)
	if err != nil {
		return fmt.Errorf("unmarshal gasstation error %v, for data %s", err, string(data))
	}
	return nil
}

// ETHGas get gasstation gas data
func (c *Client) ETHGas() (ETHGas, error) {
	var res ETHGas
	qp := "/json/ethgasAPI.json"
	err := c.doRequest(http.MethodGet, qp, &res)
	return res, err
}
