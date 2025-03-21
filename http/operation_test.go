package http

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
)

const EXPECTED_NIL_GOT = "Expected nil, got %v"
const EXPECTED_ERROR_GOT_NIL = "Expected error, got nil"
const EXPECTED_BUT_GOT = "Expected %v, got %v"
const VALID_OUTPUT = "$.output.message"

var validOutputMap = map[string]interface{}{
	"output": map[string]interface{}{
		"message": "world",
	},
}

type MockRoundTripper struct {
	MockResponse *http.Response
	MockError    error
}

// RoundTrip implements the RoundTripper interface
func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.MockResponse, m.MockError
}

func TestGetAction(t *testing.T) {
	statusCodeStr := "200"
	responsesMap := map[string]interface{}{
		"200": map[string]interface{}{
			"output": map[string]interface{}{
				"message": VALID_OUTPUT,
			},
			"next": "next",
		},
	}
	output, next, err := getActions(responsesMap, statusCodeStr)

	if err != nil {
		t.Errorf(EXPECTED_NIL_GOT, err)
	}
	if output["message"] != responsesMap[statusCodeStr].(map[string]interface{})["output"].(map[string]interface{})["message"] {
		t.Errorf(EXPECTED_BUT_GOT, VALID_OUTPUT, output)
	}
	if next != "next" {
		t.Errorf(EXPECTED_BUT_GOT, "next", next)
	}
}

func TestGetDefaultAction(t *testing.T) {
	statusCodeStr := "201"
	responsesMap := map[string]interface{}{
		"default": map[string]interface{}{
			"output": map[string]interface{}{
				"message": VALID_OUTPUT,
			},
			"next": "next",
		},
	}
	output, next, err := getActions(responsesMap, statusCodeStr)

	if err != nil {
		t.Errorf(EXPECTED_NIL_GOT, err)
	}
	if output["message"] != responsesMap["default"].(map[string]interface{})["output"].(map[string]interface{})["message"] {
		t.Errorf(EXPECTED_BUT_GOT, VALID_OUTPUT, output)
	}
	if next != "next" {
		t.Errorf(EXPECTED_BUT_GOT, "next", next)
	}
}

func TestGetActionInvalidStatusCode(t *testing.T) {
	statusCodeStr := "404"
	responsesMap := map[string]interface{}{
		"200": map[string]interface{}{
			"output": map[string]interface{}{
				"message": VALID_OUTPUT,
			},
			"next": "next",
		},
	}
	output, next, err := getActions(responsesMap, statusCodeStr)

	if err == nil {
		t.Error(EXPECTED_ERROR_GOT_NIL)
	}
	if output != nil {
		t.Errorf(EXPECTED_NIL_GOT, output)
	}
	if next != "error" {
		t.Errorf(EXPECTED_BUT_GOT, "error", next)
	}
}

func TestGetActionInvalidOutputFormat(t *testing.T) {
	statusCodeStr := "200"
	responsesMap := map[string]interface{}{
		"200": map[string]interface{}{
			"output": 1,
			"next":   "next",
		},
	}
	output, next, err := getActions(responsesMap, statusCodeStr)

	if err == nil {
		t.Error(EXPECTED_ERROR_GOT_NIL)
	}
	if output != nil {
		t.Errorf(EXPECTED_NIL_GOT, output)
	}
	if next != "error" {
		t.Errorf(EXPECTED_BUT_GOT, "error", next)
	}
}

func TestGetActionInvalidNextFormat(t *testing.T) {
	statusCodeStr := "200"
	responsesMap := map[string]interface{}{
		"200": map[string]interface{}{
			"output": map[string]interface{}{
				"message": VALID_OUTPUT,
			},
			"next": 1,
		},
	}
	output, next, err := getActions(responsesMap, statusCodeStr)

	if err == nil {
		t.Error(EXPECTED_ERROR_GOT_NIL)
	}
	if output != nil {
		t.Errorf(EXPECTED_NIL_GOT, output)
	}
	if next != "error" {
		t.Errorf(EXPECTED_BUT_GOT, "error", next)
	}
}
func TestHttpRequest(t *testing.T) {
	ctx := context.Background()
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(`{"message": "success"}`)),
		Header:     make(http.Header),
	}

	// Create a mock HTTP client
	mockTransport := &MockRoundTripper{
		MockResponse: mockResponse,
		MockError:    nil,
	}
	mockClient := &http.Client{
		Transport: mockTransport,
	}
	method := "GET"
	url := "http://example.com"
	requestBodyString := ""
	headers := map[string]interface{}{
		"Content-Type": "application/json",
	}
	stepOutputs := map[string]interface{}{}

	response, err := httpRequest(ctx, mockClient, method, url, requestBodyString, headers, stepOutputs)

	if err != nil {
		t.Errorf(EXPECTED_NIL_GOT, err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf(EXPECTED_BUT_GOT, response.StatusCode, http.StatusOK)
	}
}

func TestHttpRequestError(t *testing.T) {
	ctx := context.Background()
	mockResponse := &http.Response{}
	mockTransport := &MockRoundTripper{
		MockResponse: mockResponse,
		MockError:    http.ErrHandlerTimeout,
	}
	mockClient := &http.Client{
		Transport: mockTransport,
	}
	method := "GET"
	url := "http://example.com"
	requestBodyString := ""
	headers := map[string]interface{}{
		"Content-Type": "application/json",
	}
	stepOutputs := map[string]interface{}{}

	response, err := httpRequest(ctx, mockClient, method, url, requestBodyString, headers, stepOutputs)

	if err == nil {
		t.Error(EXPECTED_ERROR_GOT_NIL)
	}
	if response != nil {
		t.Errorf(EXPECTED_NIL_GOT, response)
	}
}

func TestHttpRequestInvalidUrl(t *testing.T) {
	ctx := context.Background()
	mockResponse := &http.Response{}
	mockTransport := &MockRoundTripper{
		MockResponse: mockResponse,
		MockError:    nil,
	}
	mockClient := &http.Client{
		Transport: mockTransport,
	}
	method := "GET"
	url := "http://example.com/$.test"
	requestBodyString := ""
	headers := map[string]interface{}{
		"Content-Type": "application/json",
	}
	stepOutputs := map[string]interface{}{}

	response, err := httpRequest(ctx, mockClient, method, url, requestBodyString, headers, stepOutputs)

	if err == nil {
		t.Error(EXPECTED_ERROR_GOT_NIL)
	}
	if response != nil {
		t.Errorf(EXPECTED_NIL_GOT, response)
	}
}

func TestHttpRequestInvalidHeader(t *testing.T) {
	ctx := context.Background()
	mockResponse := &http.Response{}
	mockTransport := &MockRoundTripper{
		MockResponse: mockResponse,
		MockError:    nil,
	}
	mockClient := &http.Client{
		Transport: mockTransport,
	}
	method := "GET"
	url := "http://example.com"
	requestBodyString := ""
	headers := map[string]interface{}{
		"Content-Type": "$.test",
	}
	stepOutputs := map[string]interface{}{}

	response, err := httpRequest(ctx, mockClient, method, url, requestBodyString, headers, stepOutputs)

	if err == nil {
		t.Error(EXPECTED_ERROR_GOT_NIL)
	}
	if response != nil {
		t.Errorf(EXPECTED_NIL_GOT, response)
	}
}

func TestRun(t *testing.T) {
	ctx := context.Background()
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(`{"message": "success"}`)),
		Header:     make(http.Header),
	}

	// Create a mock HTTP client
	mockTransport := &MockRoundTripper{
		MockResponse: mockResponse,
		MockError:    nil,
	}
	mockClient := &http.Client{
		Transport: mockTransport,
	}
	stepMap := map[string]interface{}{
		"method": "GET",
		"url":    "http://example.com",
		"body":   map[string]interface{}{},
		"headers": map[string]interface{}{
			"Content-Type": "application/json",
		},
		"responses": map[string]interface{}{
			"200": map[string]interface{}{
				"output": map[string]interface{}{
					"message": "$.body.message",
				},
				"next": "next",
			},
		},
	}
	stepOutputs := validOutputMap

	output, next, err := Run(ctx, mockClient, stepMap, stepOutputs)

	if err != nil {
		t.Errorf(EXPECTED_NIL_GOT, err)
	}
	if output.(map[string]interface{})["message"] != "success" {
		t.Errorf(EXPECTED_BUT_GOT, "success", output.(map[string]interface{})["message"])
	}
	if next != "next" {
		t.Errorf(EXPECTED_BUT_GOT, "next", next)
	}
}

func TestRunError(t *testing.T) {
	ctx := context.Background()
	mockResponse := &http.Response{}
	mockTransport := &MockRoundTripper{
		MockResponse: mockResponse,
		MockError:    http.ErrHandlerTimeout,
	}
	mockClient := &http.Client{
		Transport: mockTransport,
	}
	stepMap := map[string]interface{}{
		"method": "GET",
		"url":    "http://example.com",
		"body":   map[string]interface{}{},
		"headers": map[string]interface{}{
			"Content-Type": "application/json",
		},
		"responses": map[string]interface{}{
			"200": map[string]interface{}{
				"output": map[string]interface{}{
					"message": "$.body.message",
				},
				"next": "next",
			},
		},
	}
	stepOutputs := validOutputMap
	expectedOutput := "Get \"http://example.com\": http: Handler timeout"

	output, next, err := Run(ctx, mockClient, stepMap, stepOutputs)

	if err == nil {
		t.Error(EXPECTED_ERROR_GOT_NIL)
	}
	if output != expectedOutput {
		t.Errorf(EXPECTED_BUT_GOT, expectedOutput, output)
	}
	if next != "error" {
		t.Errorf(EXPECTED_BUT_GOT, "error", next)
	}
}

func TestRunInvalidOutputFormat(t *testing.T) {
	ctx := context.Background()
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(`{"message": "success"}`)),
		Header:     make(http.Header),
	}
	mockTransport := &MockRoundTripper{
		MockResponse: mockResponse,
		MockError:    nil,
	}
	mockClient := &http.Client{
		Transport: mockTransport,
	}
	stepMap := map[string]interface{}{
		"method": "GET",
		"url":    "http://example.com",
		"body":   map[string]interface{}{},
		"headers": map[string]interface{}{
			"Content-Type": "application/json",
		},
		"responses": map[string]interface{}{
			"200": map[string]interface{}{
				"output": 1,
				"next":   "next",
			},
		},
	}
	stepOutputs := validOutputMap
	expectedOutput := "invalid output format"

	output, next, err := Run(ctx, mockClient, stepMap, stepOutputs)

	if err == nil {
		t.Error(EXPECTED_ERROR_GOT_NIL)
	}
	if output != expectedOutput {
		t.Errorf(EXPECTED_BUT_GOT, expectedOutput, output)
	}
	if next != "error" {
		t.Errorf(EXPECTED_BUT_GOT, "error", next)
	}
}

func TestRunInvalidRequestBody(t *testing.T) {
	ctx := context.Background()
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(`{"message": "success"}`)),
		Header:     make(http.Header),
	}
	mockTransport := &MockRoundTripper{
		MockResponse: mockResponse,
		MockError:    nil,
	}
	mockClient := &http.Client{
		Transport: mockTransport,
	}
	stepMap := map[string]interface{}{
		"method": "POST",
		"url":    "http://example.com",
		"body": map[string]interface{}{
			"message": "$.body.message",
		},
		"headers": map[string]interface{}{
			"Content-Type": "application/json",
		},
		"responses": map[string]interface{}{
			"200": map[string]interface{}{
				"output": map[string]interface{}{
					"message": "$.body.message",
				},
				"next": "next",
			},
		},
	}
	stepOutputs := map[string]interface{}{
		"body": 1,
	}
	expectedOutput := "unsupported value type int for select, expected map[string]interface{} or []interface{}"

	output, next, err := Run(ctx, mockClient, stepMap, stepOutputs)

	if err == nil {
		t.Error(EXPECTED_ERROR_GOT_NIL)
	}
	if output != expectedOutput {
		t.Errorf(EXPECTED_BUT_GOT, expectedOutput, output)
	}
	if next != "error" {
		t.Errorf(EXPECTED_BUT_GOT, "error", next)
	}
}

func TestRunInvalidJsonResponse(t *testing.T) {
	ctx := context.Background()
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(`{"message": "success"`)),
		Header:     make(http.Header),
	}
	mockTransport := &MockRoundTripper{
		MockResponse: mockResponse,
		MockError:    nil,
	}
	mockClient := &http.Client{
		Transport: mockTransport,
	}
	stepMap := map[string]interface{}{
		"method": "GET",
		"url":    "http://example.com",
		"body":   map[string]interface{}{},
		"headers": map[string]interface{}{
			"Content-Type": "application/json",
		},
		"responses": map[string]interface{}{
			"200": map[string]interface{}{
				"output": map[string]interface{}{
					"message": "$.body.message",
				},
				"next": "next",
			},
		},
	}
	stepOutputs := validOutputMap
	expectedOutput := "unexpected EOF"

	output, next, err := Run(ctx, mockClient, stepMap, stepOutputs)

	if err == nil {
		t.Error(EXPECTED_ERROR_GOT_NIL)
	}
	if output != expectedOutput {
		t.Errorf(EXPECTED_BUT_GOT, expectedOutput, output)
	}
	if next != "error" {
		t.Errorf(EXPECTED_BUT_GOT, "error", next)
	}
}

func TestRunInvalidOutputMap(t *testing.T) {
	ctx := context.Background()
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(`{"message": "success"}`)),
		Header:     make(http.Header),
	}
	mockTransport := &MockRoundTripper{
		MockResponse: mockResponse,
		MockError:    nil,
	}
	mockClient := &http.Client{
		Transport: mockTransport,
	}
	stepMap := map[string]interface{}{
		"method": "GET",
		"url":    "http://example.com",
		"body":   map[string]interface{}{},
		"headers": map[string]interface{}{
			"Content-Type": "application/json",
		},
		"responses": map[string]interface{}{
			"200": map[string]interface{}{
				"output": map[string]interface{}{
					"message": "$.body.test",
				},
				"next": "next",
			},
		},
	}
	stepOutputs := map[string]interface{}{
		"output": 1,
	}
	expectedOutput := "unknown key test"

	output, next, err := Run(ctx, mockClient, stepMap, stepOutputs)

	if err == nil {
		t.Error(EXPECTED_ERROR_GOT_NIL)
	}
	if output != expectedOutput {
		t.Errorf(EXPECTED_BUT_GOT, expectedOutput, output)
	}
	if next != "error" {
		t.Errorf(EXPECTED_BUT_GOT, "error", next)
	}
}
