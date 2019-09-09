package world

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin/json"

	"github.com/KyberNetwork/reserve-data/common"
)

func (tw *TheWorld) getCoinbaseInfo(ep string) common.CoinbaseData {
	var (
		client = &http.Client{Timeout: 30 * time.Second}
		url    = ep
		result = common.CoinbaseData{}
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return common.CoinbaseData{
			Valid: false,
			Error: err.Error(),
		}
	}

	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return common.CoinbaseData{
			Valid: false,
			Error: err.Error(),
		}
	}
	defer func() {
		if cErr := resp.Body.Close(); cErr != nil {
			log.Printf("failed to close response body: %s", cErr.Error())
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return common.CoinbaseData{
			Valid: false,
			Error: fmt.Sprintf("unexpected return code: %d", resp.StatusCode),
		}
	}

	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return common.CoinbaseData{
			Valid: false,
			Error: err.Error(),
		}
	}
	result.Valid = true
	return result
}

func (tw *TheWorld) getGeminiInfo(url string) common.GeminiData {
	var (
		client = &http.Client{Timeout: 30 * time.Second}
		result = common.GeminiData{}
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return common.GeminiData{
			Valid: false,
			Error: err.Error(),
		}
	}

	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return common.GeminiData{
			Valid: false,
			Error: err.Error(),
		}
	}
	defer func() {
		if cErr := resp.Body.Close(); cErr != nil {
			log.Printf("failed to close response body: %s", cErr.Error())
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return common.GeminiData{
			Valid: false,
			Error: fmt.Sprintf("unexpected return code: %d", resp.StatusCode),
		}
	}

	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return common.GeminiData{
			Valid: false,
			Error: err.Error(),
		}
	}
	result.Valid = true
	return result
}

func (tw *TheWorld) GetBTCInfo() (common.BTCData, error) {
	return common.BTCData{
		Coinbase: tw.getCoinbaseInfo(tw.endpoint.CoinbaseBTCEndpoint()),
		Gemini:   tw.getGeminiInfo(tw.endpoint.GeminiBTCEndpoint()),
	}, nil
}
