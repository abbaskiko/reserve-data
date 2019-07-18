package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/KyberNetwork/reserve-data/common"
)

type assertFn func(t *testing.T, resp *httptest.ResponseRecorder)

type testCase struct {
	msg      string
	endpoint string
	method   string
	data     interface{}
	assert   assertFn
}

func newAssertCreated(expectedData []byte) assertFn {
	return func(t *testing.T, resp *httptest.ResponseRecorder) {
		t.Helper()

		if resp.Code != http.StatusOK {
			t.Fatalf("wrong return code, expected: %d, got: %d, body[%s]", http.StatusOK, resp.Code, resp.Body.String())
		}

		type responseBody struct {
			ID uint64
		}

		decoded := responseBody{}
		if aErr := json.NewDecoder(resp.Body).Decode(&decoded); aErr != nil {
			t.Fatal(aErr)
		}

		t.Logf("returned ID: %v", decoded.ID)
	}
}

func newAssertHTTPCode(code int) assertFn {
	return func(t *testing.T, resp *httptest.ResponseRecorder) {
		t.Helper()
		if resp.Code != code {
			t.Fatalf("wrong return code, expected: %d, got: %d, error = [%s]", code, resp.Code, resp.Body.String())
		}
		t.Logf("response: %s\n", resp.Body.String())
	}
}

func newAssertGetEquation(expectedData []byte) assertFn {
	return func(t *testing.T, resp *httptest.ResponseRecorder) {
		t.Helper()

		var expected common.PWIEquationRequestV2
		if resp.Code != http.StatusOK {
			t.Fatalf("wrong return code, expected: %d, got: %d", http.StatusOK, resp.Code)
		}

		type responseBody struct {
			Success bool
			Data    common.PWIEquationRequestV2
		}

		decoded := responseBody{}
		if aErr := json.NewDecoder(resp.Body).Decode(&decoded); aErr != nil {
			t.Fatal(aErr)
		}

		if decoded.Success != true {
			t.Errorf("wrong success status, expected: %t, got: %t", true, decoded.Success)
		}

		t.Logf("returned pending PWI equation request: %v", decoded.Data)

		if len(decoded.Data) != 2 {
			t.Fatalf("wrong number of tokens, expected: %d, got %d", 2, len(decoded.Data))
		}

		if aErr := json.Unmarshal(expectedData, &expected); aErr != nil {
			t.Fatal(aErr)
		}

		if !reflect.DeepEqual(expected, decoded.Data) {
			t.Logf("expected data doesn't match: %v, decoded: %v", expected, decoded)
		}
	}
}

func testHTTPRequest(t *testing.T, tc testCase, handler http.Handler) {
	t.Helper()

	req, tErr := http.NewRequest(tc.method, tc.endpoint, nil)
	if tErr != nil {
		t.Fatal(tErr)
	}

	data, err := json.Marshal(tc.data)
	if err != nil {
		t.Fatal(err)
	}

	if tc.data != nil {
		req.Body = ioutil.NopCloser(bytes.NewReader(data))
		req.Header.Add("Content-Type", "application/json")
	}

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	tc.assert(t, resp)
}
