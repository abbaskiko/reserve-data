package authhttp

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/KyberNetwork/httpsign-utils/sign"
	"github.com/pkg/errors"
)

var (
	hc = &http.Client{}
)

// AuthHTTP ...
type AuthHTTP struct {
	client       *http.Client
	accessKey    string
	accessSecret string
}

// NewAuthHTTP ...
func NewAuthHTTP(accessKey, accessSecret string) *AuthHTTP {
	return &AuthHTTP{
		client:       hc,
		accessKey:    accessKey,
		accessSecret: accessSecret,
	}
}

type errorResponse struct {
	Msg string `json:"msg"`
}

// DoReq do request
func (ah *AuthHTTP) DoReq(url string, method string, params map[string]string) ([]byte, error) {
	var (
		httpMethod = strings.ToUpper(method)
	)
	req, err := http.NewRequest(httpMethod, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create get request")
	}
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	if ah.accessKey != "" && ah.accessSecret != "" {
		req, err = sign.Sign(req, ah.accessKey, ah.accessSecret)
		if err != nil {
			return nil, errors.Wrap(err, "failed to sign get request")
		}
	}

	rsp, err := ah.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to do get req")
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read response's body")
	}
	if err := rsp.Body.Close(); err != nil {
		return nil, errors.Wrap(err, "cannot close response's body")
	}
	if rsp.StatusCode != 200 {
		var er errorResponse
		if err := json.Unmarshal(body, &er); err != nil {
			return nil, errors.Wrap(err, "cannot unmarshal response data")
		}
		return nil, errors.Errorf("receive unexpected code, actual code: %d, err: %s", rsp.StatusCode, er.Msg)
	}
	return body, nil
}
