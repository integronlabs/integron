package helpers

import (
	"net/http"
	"testing"
)

type mockResponseWriter struct {
	headers http.Header
}

func (m *mockResponseWriter) Header() http.Header {
	return m.headers
}

func (m *mockResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (m *mockResponseWriter) WriteHeader(int) {
	// mock function
}

func TestFillResponseHeaders(t *testing.T) {
	responseHeaders := map[string][]string{"key1": {"value1"}}
	w := &mockResponseWriter{headers: http.Header{}}
	FillResponseHeaders(responseHeaders, w)
	if w.headers.Get("key1") != "value1" {
		t.Errorf("Expected value1, got %s", w.headers.Get("key1"))
	}
}

func TestExtractParams(t *testing.T) {
	pathParams := map[string]string{"key1": "value1"}
	queryParams := map[string][]string{"key2": {"value2"}}
	params := ExtractParams(pathParams, queryParams)
	if params["key1"] != "value1" {
		t.Errorf("Expected value1, got %s", params["key1"])
	}
	if params["key2"] != "value2" {
		t.Errorf("Expected value2, got %s", params["key2"])
	}
}
