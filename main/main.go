package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	httpOperation "github.com/integronlabs/integron/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/swaggest/swgui/v5emb"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path"},
	)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}

var router routers.Router
var ctx context.Context

const INVALID_STEP_DEFINITION = "invalid step definition"

func extractParams(pathParams map[string]string, queryParams map[string][]string) map[string]interface{} {
	params := make(map[string]interface{})
	for key, value := range pathParams {
		params[key] = value
	}
	for key, value := range queryParams {
		params[key] = value[0]
	}
	return params
}

func createStepsMap(stepsArray []interface{}) (map[string]interface{}, error) {
	steps := make(map[string]interface{})
	for _, v := range stepsArray {
		stepsMap, ok := v.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf(INVALID_STEP_DEFINITION)
		}
		steps[stepsMap["name"].(string)] = stepsMap
	}
	return steps, nil
}

func fillResponseHeaders(responseHeaders http.Header, w http.ResponseWriter) {
	for k, v := range responseHeaders {
		w.Header().Set(k, v[0])
	}
}

func processStep(currentStepKey string, w http.ResponseWriter, steps map[string]interface{}, input map[string]interface{}, stepOutputs map[string]interface{}, stepInput interface{}) (interface{}, string) {
	var next string
	var err error
	step, ok := steps[currentStepKey]
	if !ok {
		return fmt.Errorf(INVALID_STEP_DEFINITION), "error"
	}
	stepMap, ok := step.(map[string]interface{})
	if !ok {
		return fmt.Errorf(INVALID_STEP_DEFINITION), "error"
	}

	switch (stepMap["type"]).(string) {
	case "http":
		stepOutputs[currentStepKey], next, err = httpOperation.Run(stepMap, input, stepOutputs)
		if err != nil {
			return err, "error"
		}
	case "error":
		message, _ := json.Marshal(map[string]interface{}{"message": stepInput})
		http.Error(w, string(message), http.StatusInternalServerError)
		return nil, "end"
	}
	return stepOutputs[currentStepKey], next
}

func handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Seconds()
		httpRequestsTotal.WithLabelValues(r.URL.Path).Inc()
		httpRequestDuration.WithLabelValues(r.URL.Path).Observe(duration)
	}()
	// Find route
	route, pathParams, err := router.FindRoute(r)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Method not found", http.StatusNotFound)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var output interface{}
	var stepInput interface{}
	stepOutputs := make(map[string]interface{})
	input := extractParams(pathParams, r.URL.Query())

	stepOutputs["request"] = input

	stepsArray, ok := route.PathItem.GetOperation(route.Method).Extensions["x-integron-steps"].([]interface{})

	if !ok {
		http.Error(w, "Invalid x-integron-steps", http.StatusInternalServerError)
		return
	}

	currentStepKey := stepsArray[0].(map[string]interface{})["name"].(string)
	steps, err := createStepsMap(stepsArray)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stepInput = input
	for {
		var next string
		stepOutputs[currentStepKey], next = processStep(currentStepKey, w, steps, input, stepOutputs, stepInput)

		if next == "" {
			output = stepOutputs[currentStepKey]
			break
		} else if next == "end" {
			return
		}
		stepInput = stepOutputs[currentStepKey]
		currentStepKey = next
	}
	responseHeaders := http.Header{
		"Content-Type":                 []string{"application/json"},
		"Access-Control-Allow-Origin":  []string{"*"},
		"Access-Control-Allow-Methods": []string{"GET, POST, PUT, DELETE"},
		"Access-Control-Allow-Headers": []string{"Content-Type"},
	}
	responseCode := 200
	jsonBody, _ := json.Marshal(output)
	fmt.Println(string(jsonBody))
	responseBody := []byte(jsonBody)

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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fillResponseHeaders(responseHeaders, w)

	w.WriteHeader(responseCode)

	w.Write(responseBody)
}

func main() {
	openapiSpec := flag.String("spec", "docs/openapi.yaml", "path to the openapi spec")

	flag.Parse()

	ctx = context.Background()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	doc, err := loader.LoadFromFile(*openapiSpec)
	if err != nil {
		panic(err)
	}

	// Validate document
	err = doc.Validate(ctx)
	if err != nil {
		panic(err)
	}

	r, err := gorillamux.NewRouter(doc)
	if err != nil {
		panic(err)
	}
	router = r

	http.Handle("/", http.HandlerFunc(handler))

	fs := http.FileServer(http.Dir("docs/"))
	http.Handle("/docs/", http.StripPrefix("/docs/", fs))

	http.Handle("/ui/", v5emb.New(
		"Integron Sunrise",
		"/"+*openapiSpec,
		"/ui/",
	))

	http.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(":8080", nil))
}
