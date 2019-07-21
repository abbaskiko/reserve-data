package http

import (
	"net/http/httptest"
	"testing"
)

type assertFn func(t *testing.T, resp *httptest.ResponseRecorder)

type testCase struct {
	msg      string
	endpoint string
	method   string
	data     map[string]string
	assert   assertFn
}
