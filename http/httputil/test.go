package httputil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type responseBody struct {
	Success bool
	Reason  string `json:"reason"`
}

// expectStatus asserts that given response matches the expected status.
func expectStatus(t *testing.T, resp *httptest.ResponseRecorder, success bool, reason string) {
	t.Helper()

	if resp.Code != http.StatusOK {
		t.Fatalf("wrong return code, expected: %d, got: %d", http.StatusOK, resp.Code)
	}
	decoded := &responseBody{}
	if err := json.NewDecoder(resp.Body).Decode(decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Success != success {
		t.Errorf("wrong success status, expected: %t, got: %t", success, decoded.Success)
	}
	if len(reason) != 0 && reason != decoded.Reason {
		t.Errorf("wrong error msg, expected: %v, got: %v", reason, decoded.Reason)
	}
}

// ExpectSuccess asserts that given response is a success response.
func ExpectSuccess(t *testing.T, resp *httptest.ResponseRecorder) {
	t.Helper()

	expectStatus(t, resp, true, "")
}

// ExpectFailure asserts that given response is a failure response.
func ExpectFailure(t *testing.T, resp *httptest.ResponseRecorder) {
	t.Helper()

	expectStatus(t, resp, false, "")
}

func ExpectFailureWithReason(reason string) func(t *testing.T, resp *httptest.ResponseRecorder) {
	return func(t *testing.T, resp *httptest.ResponseRecorder) {
		t.Helper()
		expectStatus(t, resp, false, reason)
	}
}
