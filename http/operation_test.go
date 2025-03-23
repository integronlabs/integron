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
const EXAMPLE_URL = "http://example.com"

var validOutputMap = map[string]interface{}{
	"output": map[string]interface{}{
		"message": "world",
	},
}

var validStepMap = map[string]interface{}{
	"method": "GET",
	"url":    EXAMPLE_URL,
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

func TestHttpRequestScenarios(t *testing.T) {
	tests := []struct {
		name             string
		mockResponse     *http.Response
		mockError        error
		method           string
		url              string
		requestBody      string
		headers          map[string]interface{}
		stepOutputs      map[string]interface{}
		expectedError    bool
		expectedResponse *http.Response
	}{
		{
			name: "Valid Request",
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"message": "success"}`)),
				Header:     make(http.Header),
			},
			mockError:        nil,
			method:           "GET",
			url:              EXAMPLE_URL,
			requestBody:      "",
			headers:          map[string]interface{}{"Content-Type": "application/json"},
			stepOutputs:      map[string]interface{}{},
			expectedError:    false,
			expectedResponse: &http.Response{StatusCode: http.StatusOK},
		},
		{
			name:             "Request Timeout",
			mockResponse:     &http.Response{},
			mockError:        http.ErrHandlerTimeout,
			method:           "GET",
			url:              EXAMPLE_URL,
			requestBody:      "",
			headers:          map[string]interface{}{"Content-Type": "application/json"},
			stepOutputs:      map[string]interface{}{},
			expectedError:    true,
			expectedResponse: nil,
		},
		{
			name:             "Invalid URL",
			mockResponse:     &http.Response{},
			mockError:        nil,
			method:           "GET",
			url:              "http://example.com/$.test",
			requestBody:      "",
			headers:          map[string]interface{}{"Content-Type": "application/json"},
			stepOutputs:      map[string]interface{}{},
			expectedError:    true,
			expectedResponse: nil,
		},
		{
			name:             "Invalid Header",
			mockResponse:     &http.Response{},
			mockError:        nil,
			method:           "GET",
			url:              EXAMPLE_URL,
			requestBody:      "",
			headers:          map[string]interface{}{"Content-Type": "$.test"},
			stepOutputs:      map[string]interface{}{},
			expectedError:    true,
			expectedResponse: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockTransport := &MockRoundTripper{
				MockResponse: test.mockResponse,
				MockError:    test.mockError,
			}
			mockClient := &http.Client{
				Transport: mockTransport,
			}

			response, err := httpRequest(context.Background(), mockClient, test.method, test.url, test.requestBody, test.headers, test.stepOutputs)

			if test.expectedError && err == nil {
				t.Error(EXPECTED_ERROR_GOT_NIL)
			}
			if !test.expectedError && err != nil {
				t.Errorf(EXPECTED_NIL_GOT, err)
			}
			if test.expectedResponse != nil && response != nil && response.StatusCode != test.expectedResponse.StatusCode {
				t.Errorf(EXPECTED_BUT_GOT, test.expectedResponse.StatusCode, response.StatusCode)
			}
			if test.expectedResponse == nil && response != nil {
				t.Errorf(EXPECTED_NIL_GOT, response)
			}
		})
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
	url := EXAMPLE_URL
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
	stepMap := validStepMap
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
	stepMap := validStepMap
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
	stepMap := validStepMap
	stepMap["responses"].(map[string]interface{})["200"].(map[string]interface{})["output"] = 1

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
		"url":    EXAMPLE_URL,
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
	stepMap := validStepMap
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
	stepMap := validStepMap
	stepMap["responses"].(map[string]interface{})["200"].(map[string]interface{})["output"] = map[string]interface{}{
		"message": "$.body.test",
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
