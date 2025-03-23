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
const MESSAGE_JSON_PATH = "$.body.message"
const HEADER_CONTENT_TYPE = "Content-Type"

var validOutputMap = map[string]interface{}{
	"output": map[string]interface{}{
		"message": "world",
	},
}

var headersMap = map[string]interface{}{
	HEADER_CONTENT_TYPE: "application/json",
}

var validStepMap = map[string]interface{}{
	"method":  "GET",
	"url":     EXAMPLE_URL,
	"body":    map[string]interface{}{},
	"headers": headersMap,
	"responses": map[string]interface{}{
		"200": map[string]interface{}{
			"output": map[string]interface{}{
				"message": MESSAGE_JSON_PATH,
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
			headers:          headersMap,
			stepOutputs:      map[string]interface{}{},
			expectedError:    false,
			expectedResponse: &http.Response{StatusCode: http.StatusOK},
		},
		{
			name:             "Request Timeout",
			mockResponse:     nil,
			mockError:        http.ErrHandlerTimeout,
			method:           "GET",
			url:              EXAMPLE_URL,
			requestBody:      "",
			headers:          headersMap,
			stepOutputs:      map[string]interface{}{},
			expectedError:    true,
			expectedResponse: nil,
		},
		{
			name:             "Invalid URL",
			mockResponse:     nil,
			mockError:        nil,
			method:           "GET",
			url:              "http://example.com/$.test",
			requestBody:      "",
			headers:          headersMap,
			stepOutputs:      map[string]interface{}{},
			expectedError:    true,
			expectedResponse: nil,
		},
		{
			name:             "Invalid Header",
			mockResponse:     nil,
			mockError:        nil,
			method:           "GET",
			url:              EXAMPLE_URL,
			requestBody:      "",
			headers:          map[string]interface{}{HEADER_CONTENT_TYPE: "$.test"},
			stepOutputs:      map[string]interface{}{},
			expectedError:    true,
			expectedResponse: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockClient := createMockClient(test.mockResponse, test.mockError)

			response, err := httpRequest(context.Background(), mockClient, test.method, test.url, test.requestBody, test.headers, test.stepOutputs)

			assertError(t, test.expectedError, err)
			assertResponse(t, test.expectedResponse, response)
		})
	}
}

func createMockClient(mockResponse *http.Response, mockError error) *http.Client {
	return &http.Client{
		Transport: &MockRoundTripper{
			MockResponse: mockResponse,
			MockError:    mockError,
		},
	}
}

func assertError(t *testing.T, expectedError bool, err error) {
	if expectedError && err == nil {
		t.Error(EXPECTED_ERROR_GOT_NIL)
	}
	if !expectedError && err != nil {
		t.Errorf(EXPECTED_NIL_GOT, err)
	}
}

func assertResponse(t *testing.T, expectedResponse, actualResponse *http.Response) {
	if expectedResponse != nil && actualResponse != nil {
		if expectedResponse.StatusCode != actualResponse.StatusCode {
			t.Errorf(EXPECTED_BUT_GOT, expectedResponse.StatusCode, actualResponse.StatusCode)
		}
	}
	if expectedResponse == nil && actualResponse != nil {
		t.Errorf(EXPECTED_NIL_GOT, actualResponse)
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
