package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/integronlabs/integron/helpers"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_worker_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"url"},
	)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_worker_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"url"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}

func getActions(stepMap map[string]interface{}, statusCodeStr string) (map[string]interface{}, string, error) {
	statusMap, ok := stepMap[statusCodeStr].(map[string]interface{})
	if !ok {
		return map[string]interface{}{}, "error", fmt.Errorf("could not find actions for status %s", statusCodeStr)
	}
	outputMap, ok := statusMap["output"].(map[string]interface{})
	if !ok {
		return map[string]interface{}{}, "error", fmt.Errorf("invalid output format")
	}
	next, ok := stepMap[statusCodeStr].(map[string]interface{})["next"].(string)
	if !ok {
		return map[string]interface{}{}, "error", fmt.Errorf("invalid next format")
	}
	return outputMap, next, nil
}

func httpRequest(ctx context.Context, client *http.Client, method string, url string, requestBodyString string, headers map[string]interface{}, stepOutputs map[string]interface{}) (*http.Response, error) {
	url = helpers.Replace(url, stepOutputs)

	httpRequest, err := http.NewRequestWithContext(ctx, method, url, strings.NewReader(requestBodyString))
	if err != nil {
		return nil, err
	}
	// set headers
	for key, value := range headers {
		value := helpers.Replace(value.(string), stepOutputs)
		httpRequest.Header.Set(key, value)
	}
	start := time.Now()
	response, err := client.Do(httpRequest)

	if err != nil {
		return nil, err
	}

	defer func() {
		duration := time.Since(start).Seconds()
		httpRequestsTotal.WithLabelValues(httpRequest.URL.Path).Inc()
		httpRequestDuration.WithLabelValues(httpRequest.URL.Path).Observe(duration)
	}()

	return response, nil
}

func Run(ctx context.Context, client *http.Client, stepMap map[string]interface{}, stepOutputs map[string]interface{}) (interface{}, string, error) {
	// get values
	method, _ := stepMap["method"].(string)
	url, _ := stepMap["url"].(string)
	requestBodyMap, _ := stepMap["body"].(map[string]interface{})
	headers, _ := stepMap["headers"].(map[string]interface{})

	requestBody, err := helpers.TransformBody(stepOutputs, requestBodyMap)
	if err != nil {
		return err.Error(), "error", nil
	}
	requestBodyJson, _ := json.Marshal(requestBody)
	requestBodyString := string(requestBodyJson)

	response, err := httpRequest(ctx, client, method, url, requestBodyString, headers, stepOutputs)

	if err != nil {
		return err.Error(), "error", nil
	}

	defer response.Body.Close()

	var responseData interface{}
	if err := json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		return err.Error(), "error", nil
	}

	responseMap := map[string]interface{}{
		"status":  response.StatusCode,
		"headers": response.Header,
		"body":    responseData,
	}

	statusCodeStr := fmt.Sprintf("%d", response.StatusCode)

	outputMap, next, err := getActions(stepMap, statusCodeStr)

	if err != nil {
		return err.Error(), "error", nil
	}

	body, err := helpers.TransformBody(responseMap, outputMap)

	if err != nil {
		return err.Error(), "error", nil
	}

	return body, next, nil
}
