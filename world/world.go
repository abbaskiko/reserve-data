package world

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/KyberNetwork/reserve-data/common"
)

//TheWorld is the concrete implementation of fetcher.TheWorld interface.
type TheWorld struct {
	endpoint Endpoint
}

func (tw *TheWorld) getBinanceInfo(ep string) common.BinanceData {
	var (
		client = &http.Client{Timeout: 30 * time.Second}
		url    = ep
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

func (tw *TheWorld) getOneForgeGoldUSDInfo() common.OneForgeGoldData {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	url := tw.endpoint.OneForgeGoldUSDDataEndpoint()
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	var err error
	var respBody []byte
	log.Printf("request to gold feed endpoint: %s", req.URL)
	resp, err := client.Do(req)
	if err != nil {
		return common.OneForgeGoldData{
			Error:   true,
			Message: err.Error(),
		}
	}
	defer func() {
		if cErr := resp.Body.Close(); cErr != nil {
			log.Printf("Close http body error: %s", cErr.Error())
		}
	}()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("gold feed returned with code: %d", resp.StatusCode)
		return common.OneForgeGoldData{
			Error:   true,
			Message: err.Error(),
		}
	}
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return common.OneForgeGoldData{
			Error:   true,
			Message: err.Error(),
		}
	}
	log.Printf("request to %s, got response from gold feed %s", req.URL, respBody)
	result := common.OneForgeGoldData{}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		result.Error = true
		result.Message = err.Error()
	}
	return result
}

func (tw *TheWorld) getOneForgeGoldETHInfo() common.OneForgeGoldData {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	url := tw.endpoint.OneForgeGoldETHDataEndpoint()
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	var err error
	var respBody []byte
	log.Printf("request to gold feed endpoint: %s", req.URL)
	resp, err := client.Do(req)
	if err != nil {
		return common.OneForgeGoldData{
			Error:   true,
			Message: err.Error(),
		}
	}
	defer func() {
		if cErr := resp.Body.Close(); cErr != nil {
			log.Printf("Response body close error: %s", cErr.Error())
		}
	}()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("gold feed returned with code: %d", resp.StatusCode)
		return common.OneForgeGoldData{
			Error:   true,
			Message: err.Error(),
		}
	}
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return common.OneForgeGoldData{
			Error:   true,
			Message: err.Error(),
		}
	}
	log.Printf("request to %s, got response from gold feed %s", req.URL, respBody)
	result := common.OneForgeGoldData{}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		result.Error = true
		result.Message = err.Error()
	}
	return result
}

func (tw *TheWorld) getDGXGoldInfo() common.DGXGoldData {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	url := tw.endpoint.GoldDataEndpoint()
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	var err error
	var respBody []byte
	log.Printf("request to gold feed endpoint: %s", req.URL)
	resp, err := client.Do(req)
	if err != nil {
		return common.DGXGoldData{
			Valid: false,
			Error: err.Error(),
		}
	}
	defer func() {
		if cErr := resp.Body.Close(); cErr != nil {
			log.Printf("Close reponse body error: %s", cErr.Error())
		}
	}()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("gold feed returned with code: %d", resp.StatusCode)
		return common.DGXGoldData{
			Valid: false,
			Error: err.Error(),
		}
	}
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return common.DGXGoldData{
			Valid: false,
			Error: err.Error(),
		}
	}
	log.Printf("request to %s, got response from gold feed %s", req.URL, respBody)
	result := common.DGXGoldData{
		Valid: true,
	}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		result.Valid = false
		result.Error = err.Error()
	}
	return result
}

func (tw *TheWorld) getGDAXGoldInfo() common.GDAXGoldData {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	url := tw.endpoint.GDAXDataEndpoint()
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	var err error
	var respBody []byte
	log.Printf("request to gold feed endpoint: %s", req.URL)
	resp, err := client.Do(req)
	if err != nil {
		return common.GDAXGoldData{
			Valid: false,
			Error: err.Error(),
		}
	}
	defer func() {
		if cErr := resp.Body.Close(); cErr != nil {
			log.Printf("Response body close error: %s", cErr.Error())
		}
	}()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("gold feed returned with code: %d", resp.StatusCode)
		return common.GDAXGoldData{
			Valid: false,
			Error: err.Error(),
		}
	}
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return common.GDAXGoldData{
			Valid: false,
			Error: err.Error(),
		}
	}
	log.Printf("request to %s, got response from gold feed %s", req.URL, respBody)
	result := common.GDAXGoldData{
		Valid: true,
	}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		result.Valid = false
		result.Error = err.Error()
	}
	return result
}

func (tw *TheWorld) getKrakenGoldInfo() common.KrakenGoldData {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	url := tw.endpoint.KrakenDataEndpoint()
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	var err error
	var respBody []byte
	log.Printf("request to gold feed endpoint: %s", req.URL)
	resp, err := client.Do(req)
	if err != nil {
		return common.KrakenGoldData{
			Valid: false,
			Error: err.Error(),
		}
	}
	defer func() {
		if cErr := resp.Body.Close(); cErr != nil {
			log.Printf("Response body close error: %s", cErr.Error())
		}
	}()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("gold feed returned with code: %d", resp.StatusCode)
		return common.KrakenGoldData{
			Valid: false,
			Error: err.Error(),
		}
	}
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return common.KrakenGoldData{
			Valid: false,
			Error: err.Error(),
		}
	}
	log.Printf("request to %s, got response from gold feed %s", req.URL, respBody)
	result := common.KrakenGoldData{
		Valid: true,
	}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		result.Valid = false
		result.Error = err.Error()
	}
	return result
}

func (tw *TheWorld) getGeminiGoldInfo() common.GeminiGoldData {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	url := tw.endpoint.GeminiDataEndpoint()
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	var err error
	var respBody []byte
	log.Printf("request to gold feed endpoint: %s", req.URL)
	resp, err := client.Do(req)
	if err != nil {
		return common.GeminiGoldData{
			Valid: false,
			Error: err.Error(),
		}
	}
	defer func() {
		if cErr := resp.Body.Close(); cErr != nil {
			log.Printf("Response body close error: %s", cErr.Error())
		}
	}()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("gold feed returned with code: %d", resp.StatusCode)
		return common.GeminiGoldData{
			Valid: false,
			Error: err.Error(),
		}
	}
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return common.GeminiGoldData{
			Valid: false,
			Error: err.Error(),
		}
	}
	log.Printf("request to %s, got response from gold feed %s", req.URL, respBody)
	result := common.GeminiGoldData{
		Valid: true,
	}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		result.Valid = false
		result.Error = err.Error()
	}
	return result
}

func (tw *TheWorld) GetGoldInfo() (common.GoldData, error) {
	return common.GoldData{
		DGX:         tw.getDGXGoldInfo(),
		OneForgeETH: tw.getOneForgeGoldETHInfo(),
		OneForgeUSD: tw.getOneForgeGoldUSDInfo(),
		GDAX:        tw.getGDAXGoldInfo(),
		Kraken:      tw.getKrakenGoldInfo(),
		Gemini:      tw.getGeminiGoldInfo(),
	}, nil
}

//NewTheWorld return new world instance
func NewTheWorld(env string, keyfile string) (*TheWorld, error) {
	switch env {
	case common.DevMode, common.KovanMode, common.MainnetMode, common.ProductionMode, common.StagingMode, common.RopstenMode, common.AnalyticDevMode:
		endpoint, err := NewRealEndpointFromFile(keyfile)
		if err != nil {
			return nil, err
		}
		return &TheWorld{endpoint}, nil
	case common.SimulationMode:
		return &TheWorld{SimulatedEndpoint{}}, nil
	}
	panic("unsupported environment")
}
