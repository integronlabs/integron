package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/yalp/jsonpath"

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

func Run(stepMap map[string]interface{}, input map[string]interface{}, stepOutputs map[string]interface{}) (interface{}, string, error) {
	start := time.Now()

	// get values
	method, _ := stepMap["method"].(string)
	url, _ := stepMap["url"].(string)
	requestBody, _ := stepMap["body"].(map[string]interface{})
	requestBodyJson, _ := json.Marshal(requestBody)
	requestBodyString := string(requestBodyJson)
	headers, _ := stepMap["headers"].(map[string]interface{})

	body := make(map[string]interface{})

	url, err := helpers.Replace(url, stepOutputs)
	if err != nil {
		return err.Error(), "error", nil
	}
	requestBodyString, err = helpers.Replace(requestBodyString, stepOutputs)
	if err != nil {
		return err.Error(), "error", nil
	}

	fmt.Println(url, requestBodyString)

	client := &http.Client{}

	httpRequest, err := http.NewRequest(method, url, strings.NewReader(requestBodyString))
	defer func() {
		duration := time.Since(start).Seconds()
		httpRequestsTotal.WithLabelValues(httpRequest.URL.Path).Inc()
		httpRequestDuration.WithLabelValues(httpRequest.URL.Path).Observe(duration)
	}()
	if err != nil {
		return err.Error(), "error", nil
	}
	// set headers
	for key, value := range headers {
		value, err := helpers.Replace(value.(string), stepOutputs)
		if err != nil {
			return err.Error(), "error", nil
		}
		httpRequest.Header.Set(key, value)
	}
	response, err := client.Do(httpRequest)

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
	statusMap, ok := stepMap[statusCodeStr].(map[string]interface{})
	if !ok {
		return "Invalid status format", "error", nil
	}
	outputMap, ok := statusMap["output"].(map[string]interface{})
	if !ok {
		return "Invalid output format", "error", nil
	}
	next, ok := stepMap[statusCodeStr].(map[string]interface{})["next"].(string)
	if !ok {
		return "Invalid next format", "error", nil
	}

	for key, path := range outputMap {
		if value, err := jsonpath.Read(responseMap, path.(string)); err == nil {
			body[key] = value
		}
		if err != nil {
			return err.Error(), "error", nil
		}
	}

	return body, next, nil
}
