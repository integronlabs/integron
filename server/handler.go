package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/integronlabs/integron/helpers"
	"github.com/sirupsen/logrus"
)

func getStatusCode(statusCodeInterface interface{}) int {
	if status, ok := statusCodeInterface.(int); ok {
		// conver status to int
		return status
	}
	if status, ok := statusCodeInterface.(string); ok {
		// conver status to int
		statusCode, err := strconv.Atoi(status)
		if err != nil {
			return 200
		}
		return statusCode
	}
	if status, ok := statusCodeInterface.(float64); ok {
		// conver status to int
		return int(status)
	}
	return 200
}

func Error(w http.ResponseWriter, message string, code int) {
	h := w.Header()

	h.Del("Content-Length")

	h.Set("X-Content-Type-Options", "nosniff")

	responseHeaders := http.Header{
		"Content-Type":                 []string{"application/json"},
		"Access-Control-Allow-Origin":  []string{"*"},
		"Access-Control-Allow-Methods": []string{"GET, POST, PUT, DELETE"},
		"Access-Control-Allow-Headers": []string{"Content-Type"},
	}

	helpers.FillResponseHeaders(responseHeaders, w)

	body := map[string]interface{}{
		"message": message,
	}

	jsonBody, _ := json.Marshal(body)
	responseBody := []byte(jsonBody)

	w.WriteHeader(code)

	w.Write(responseBody)

	logrus.Errorf("Error: %s", message)
	logrus.Errorf("Code: %d", code)
}

func (s *Server) Handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	s.LogHandler(ctx, r)

	// Find route
	route, pathParams, err := s.Router.FindRoute(r)
	if err != nil {
		Error(w, "Method not found", http.StatusNotFound)
		return
	}

	// Validate request
	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    r,
		PathParams: pathParams,
		Route:      route,
	}

	err = openapi3filter.ValidateRequest(ctx, requestValidationInput)

	if err != nil {
		Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var output interface{}
	var stepInput interface{}
	stepOutputs := make(map[string]interface{})
	input := helpers.ExtractParams(pathParams, r.URL.Query())

	_ = json.NewDecoder(r.Body).Decode(&input)

	stepOutputs["request"] = input

	stepsArray, ok := route.PathItem.GetOperation(route.Method).Extensions["x-integron-steps"].([]interface{})

	if !ok {
		Error(w, "Invalid x-integron-steps", http.StatusInternalServerError)
		return
	}

	currentStepKey := stepsArray[0].(map[string]interface{})["name"].(string)
	steps, err := helpers.CreateStepsMap(stepsArray)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stepInput = input
	for {
		var next string
		stepOutputs[currentStepKey], next = s.ProcessStep(ctx, currentStepKey, w, steps, stepOutputs, stepInput)

		if next == "" {
			output = stepOutputs[currentStepKey]
			break
		} else if next == "end" {
			return
		}
		stepInput = stepOutputs[currentStepKey]
		currentStepKey = next
	}

	outputMap, ok := output.(map[string]interface{})
	if !ok {
		Error(w, "Invalid output format", http.StatusInternalServerError)
		return
	}
	responseCode := getStatusCode(outputMap["status"])
	jsonBody, _ := json.Marshal(outputMap["body"])
	responseBody := []byte(jsonBody)

	responseHeaders := http.Header{
		"Content-Type":                 []string{"application/json"},
		"Access-Control-Allow-Origin":  []string{"*"},
		"Access-Control-Allow-Methods": []string{"GET, POST, PUT, DELETE"},
		"Access-Control-Allow-Headers": []string{"Content-Type"},
	}
	if headers, ok := outputMap["headers"].(map[string]interface{}); ok {
		for key, value := range headers {
			responseHeaders.Set(key, value.(string))
		}
	}

	// Validate response
	responseValidationInput := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: requestValidationInput,
		Status:                 responseCode,
		Header:                 responseHeaders,
		// Body:                   io.NopCloser(strings.NewReader(string(responseBody))),
	}
	responseValidationInput.SetBodyBytes(responseBody)
	err = openapi3filter.ValidateResponse(ctx, responseValidationInput)

	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helpers.FillResponseHeaders(responseHeaders, w)

	w.WriteHeader(responseCode)

	w.Write(responseBody)
}
