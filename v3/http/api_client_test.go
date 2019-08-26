package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/KyberNetwork/reserve-data/v3/common"
)

type apiClient struct {
	s *Server
}

type apiStatus struct {
	ID      uint64          `json:"id"`
	Success bool            `json:"success"`
	Reason  string          `json:"reason"`
	Data    json.RawMessage `json:"data"`
}

type assetResponse struct {
	apiStatus
	asset common.Asset
}

type assetsResponse struct {
	apiStatus
	assets []common.Asset
}

func (c *apiClient) getAsset(id uint64) (assetResponse, error) {
	req, err := createRequest(http.MethodGet, fmt.Sprintf("/v3/asset/%d", id), nil)
	if err != nil {
		return assetResponse{}, err
	}
	httpResp := httptest.NewRecorder()
	c.s.r.ServeHTTP(httpResp, req)
	if httpResp.Code != http.StatusOK {
		return assetResponse{}, fmt.Errorf("server return %d - %s", httpResp.Code, httpResp.Body.String())
	}
	asset := common.Asset{}
	status, err := readResponse(httpResp.Body, &asset)
	if err != nil {
		return assetResponse{}, err
	}
	return assetResponse{
		apiStatus: status,
		asset:     asset,
	}, nil
}

func (c *apiClient) getAssets() (assetsResponse, error) {
	req, err := createRequest(http.MethodGet, fmt.Sprintf("/v3/asset"), nil)
	if err != nil {
		return assetsResponse{}, err
	}
	httpResp := httptest.NewRecorder()
	c.s.r.ServeHTTP(httpResp, req)
	if httpResp.Code != http.StatusOK {
		return assetsResponse{}, fmt.Errorf("server return %d - %s", httpResp.Code, httpResp.Body.String())
	}
	var assets []common.Asset
	status, err := readResponse(httpResp.Body, &assets)
	if err != nil {
		return assetsResponse{}, err
	}
	return assetsResponse{
		apiStatus: status,
		assets:    assets,
	}, nil
}

func (c *apiClient) createSettingChange(change common.SettingChange) (apiStatus, error) {
	req, err := createRequest(http.MethodPost, settingChangePath, change)
	if err != nil {
		return apiStatus{}, err
	}
	resp := httptest.NewRecorder()
	c.s.r.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		return apiStatus{}, fmt.Errorf("server return %d - %s", resp.Code, resp.Body.String())
	}
	return readResponse(resp.Body, nil)
}

func (c *apiClient) confirmSettingChange(id uint64) (apiStatus, error) {
	req, err := createRequest(http.MethodPut, fmt.Sprintf("%s/%d", settingChangePath, id), nil)
	if err != nil {
		return apiStatus{}, err
	}
	resp := httptest.NewRecorder()
	c.s.r.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		return apiStatus{}, fmt.Errorf("server return %d - %s", resp.Code, resp.Body.String())
	}
	return readResponse(resp.Body, nil)
}

// ReadResponse retry to parse response into object
func readResponse(data io.Reader, dataField interface{}) (apiStatus, error) {
	resp := apiStatus{}
	err := json.NewDecoder(data).Decode(&resp)
	if err != nil {
		return apiStatus{}, err
	}
	if dataField != nil {
		err = json.Unmarshal(resp.Data, dataField)
		if err != nil {
			return apiStatus{}, err
		}
	}
	return resp, nil
}

func createRequest(method, url string, data interface{}) (*http.Request, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return http.NewRequest(method, url, bytes.NewBuffer(body))
}
