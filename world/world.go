package world

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/KyberNetwork/reserve-data/common"
)

// TheWorld is the concrete implementation of fetcher.TheWorld interface.
type TheWorld struct {
	endpoint common.WorldEndpoints
	l        *zap.SugaredLogger
}

func (tw *TheWorld) getPublic(url string, dst interface{}) error {
	var (
		client = &http.Client{Timeout: 30 * time.Second}
	)
	caller := common.GetCallerFunctionName()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	tw.l.Debugf("%s fetch %s", caller, url)

	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		tw.l.Errorw("request on failed", "caller", caller, "err", err)
		return err
	}
	defer func() {
		if cErr := resp.Body.Close(); cErr != nil {
			tw.l.Warnw("failed to close response body", "err", cErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected return code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		tw.l.Errorw("failed to read body", "caller", caller, "err", err)
		return fmt.Errorf("failed to read response body")
	}

	if err = json.NewDecoder(bytes.NewBuffer(body)).Decode(dst); err != nil {
		tw.l.Errorw("decode failed", "caller", caller, "err", err, "body", string(common.TruncStr(body)))
		return err
	}

	return nil
}

func (tw *TheWorld) getOneForgeGoldUSDInfo() common.OneForgeGoldData {
	url := tw.endpoint.OneForgeGoldUSD.URL
	result := common.OneForgeGoldData{}
	err := tw.getPublic(url, &result)
	if err != nil {
		result.Error = true
		result.Message = err.Error()
	}
	return result
}

func (tw *TheWorld) getOneForgeGoldETHInfo() common.OneForgeGoldData {
	url := tw.endpoint.OneForgeGoldETH.URL
	result := common.OneForgeGoldData{}
	err := tw.getPublic(url, &result)
	if err != nil {
		result.Error = true
		result.Message = err.Error()
	}
	return result
}

func (tw *TheWorld) getGDAXGoldInfo() common.GDAXGoldData {
	url := tw.endpoint.GDAXData.URL
	result := common.GDAXGoldData{
		Valid: true,
	}
	err := tw.getPublic(url, &result)
	if err != nil {
		result.Valid = false
		result.Error = err.Error()
	}
	return result
}

func (tw *TheWorld) getKrakenGoldInfo() common.KrakenGoldData {
	url := tw.endpoint.KrakenData.URL
	result := common.KrakenGoldData{
		Valid: true,
	}
	err := tw.getPublic(url, &result)
	if err != nil {
		result.Valid = false
		result.Error = err.Error()
	}
	return result
}

func (tw *TheWorld) getGeminiGoldInfo() common.GeminiGoldData {

	url := tw.endpoint.GeminiData.URL
	result := common.GeminiGoldData{
		Valid: true,
	}
	err := tw.getPublic(url, &result)
	if err != nil {
		result.Valid = false
		result.Error = err.Error()
	}
	return result
}

func (tw *TheWorld) GetGoldInfo() (common.GoldData, error) {
	return common.GoldData{
		OneForgeETH: tw.getOneForgeGoldETHInfo(),
		OneForgeUSD: tw.getOneForgeGoldUSDInfo(),
		GDAX:        tw.getGDAXGoldInfo(),
		Kraken:      tw.getKrakenGoldInfo(),
		Gemini:      tw.getGeminiGoldInfo(),
	}, nil
}

func (tw *TheWorld) getFeedProviderInfo(ep string) common.FeedProviderResponse {
	var (
		url    = ep
		result = common.FeedProviderResponse{}
	)
	if err := tw.getPublic(url, &result); err != nil {
		result.Valid = false
		result.Error = err.Error()
	}
	return result
}

// NewTheWorld return new world instance
func NewTheWorld(worldEndpoints common.WorldEndpoints) *TheWorld {
	return &TheWorld{
		endpoint: worldEndpoints,
		l:        zap.S(),
	}
}
