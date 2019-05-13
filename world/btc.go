package world

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin/json"

	"github.com/KyberNetwork/reserve-data/common"
)

func (tw *TheWorld) getCoinbaseInfo() common.CoinbaseData {
	var (
		client = &http.Client{Timeout: 30 * time.Second}
		url    = tw.endpoint.CoinbaseBTCEndpoint()
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

func (tw *TheWorld) getGeminiInfo() common.GeminiData {
	var (
		client = &http.Client{Timeout: 30 * time.Second}
		url    = tw.endpoint.GeminiBTCEndpoint()
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

func (tw *TheWorld) getBitfinexInfo() common.BitfinexData {
	var (
		client = &http.Client{Timeout: 30 * time.Second}
		url    = tw.endpoint.BitfinexBTCEndpoint()
		result = common.BitfinexData{}
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return common.BitfinexData{
			Valid: false,
			Error: err.Error(),
		}
	}
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return common.BitfinexData{
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
		return common.BitfinexData{
			Valid: false,
			Error: fmt.Sprintf("unexpected return code: %d", resp.StatusCode),
		}
	}

	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return common.BitfinexData{
			Valid: false,
			Error: err.Error(),
		}
	}
	result.Valid = true
	return result
}

func (tw *TheWorld) getBinanceInfo() common.BinanceData {
	var (
		client = &http.Client{Timeout: 30 * time.Second}
		url    = tw.endpoint.BinanceBTCEndpoint()
		result = common.BinanceData{}
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return common.BinanceData{
			Valid: false,
			Error: err.Error(),
		}
	}
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return common.BinanceData{
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
		return common.BinanceData{
			Valid: false,
			Error: fmt.Sprintf("unexpected return code: %d", resp.StatusCode),
		}
	}

	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return common.BinanceData{
			Valid: false,
			Error: err.Error(),
		}
	}
	result.Valid = true
	return result
}

func (tw *TheWorld) GetBTCInfo() (common.BTCData, error) {
	return common.BTCData{
		Bitfinex: tw.getBitfinexInfo(),
		Binance:  tw.getBinanceInfo(),
		Coinbase: tw.getCoinbaseInfo(),
		Gemini:   tw.getGeminiInfo(),
	}, nil
}
