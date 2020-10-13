package authhttp

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
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

// DoReq do request
func (ah *AuthHTTP) DoReq(url string, method string, params map[string]string) ([]byte, error) {
	var (
		httpMethod = strings.ToUpper(method)
		body       io.Reader
	)
	if method != http.MethodGet {
		data, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}
		if len(data) > 0 {
			body = bytes.NewBuffer(data)
		}
	}
	req, err := http.NewRequest(httpMethod, url, body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create get request")
	}

	switch httpMethod {
	case http.MethodPost, http.MethodPut:
		req.Header.Add("Content-Type", "application/json")
	case http.MethodGet:
		q := req.URL.Query()
		for k, v := range params {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	default:
		return nil, errors.Errorf("invalid method %s", httpMethod)
	}

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

	defer func() {
		if cErr := rsp.Body.Close(); cErr != nil {
			log.Printf("failed to close body: %s", cErr.Error())
		}
	}()

	if rsp.StatusCode != 200 {
		return nil, errors.Errorf("receive unexpected code, actual code: %d", rsp.StatusCode)
	}

	return ioutil.ReadAll(rsp.Body)
}
