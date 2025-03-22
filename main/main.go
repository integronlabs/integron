package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/integronlabs/integron/helpers"

	arrayOperation "github.com/integronlabs/integron/array"
	httpOperation "github.com/integronlabs/integron/http"
	objectOperation "github.com/integronlabs/integron/object"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/swaggest/swgui/v5emb"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/sirupsen/logrus"
)

var router routers.Router
var ctx context.Context

func processStep(currentStepKey string, w http.ResponseWriter, steps map[string]interface{}, stepOutputs map[string]interface{}, stepInput interface{}) (interface{}, string) {
	logrus.Infof("Processing step: %s", currentStepKey)
	var next string
	var err error
	step, ok := steps[currentStepKey]
	if !ok {
		return fmt.Errorf(helpers.INVALID_STEP_DEFINITION), "error"
	}
	stepMap, ok := step.(map[string]interface{})
	if !ok {
		return fmt.Errorf(helpers.INVALID_STEP_DEFINITION), "error"
	}

	var stepOutput interface{}

	switch (stepMap["type"]).(string) {
	case "http":
		client := http.Client{}
		stepOutput, next, err = httpOperation.Run(ctx, &client, stepMap, stepOutputs)
		if err != nil {
			return err.Error(), "error"
		}
	case "array":
		stepOutput, next, err = arrayOperation.Run(ctx, stepMap, stepOutputs)
		if err != nil {
			return err.Error(), "error"
		}
	case "object":
		stepOutput, next, err = objectOperation.Run(ctx, stepMap, stepOutputs)
		if err != nil {
			return err.Error(), "error"
		}
	case "error":
		message, _ := json.Marshal(map[string]interface{}{"message": stepInput})
		http.Error(w, string(message), http.StatusInternalServerError)
		return nil, "end"
	}
	logrus.Infof("Step %s completed", currentStepKey)
	logrus.Infof("Step outputs: %v", stepOutputs[currentStepKey])
	return stepOutput, next
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Find route
	route, pathParams, err := router.FindRoute(r)
	if err != nil {
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
	input := helpers.ExtractParams(pathParams, r.URL.Query())

	stepOutputs["request"] = input

	stepsArray, ok := route.PathItem.GetOperation(route.Method).Extensions["x-integron-steps"].([]interface{})

	if !ok {
		http.Error(w, "Invalid x-integron-steps", http.StatusInternalServerError)
		return
	}

	currentStepKey := stepsArray[0].(map[string]interface{})["name"].(string)
	steps, err := helpers.CreateStepsMap(stepsArray)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stepInput = input
	for {
		var next string
		stepOutputs[currentStepKey], next = processStep(currentStepKey, w, steps, stepOutputs, stepInput)

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
		http.Error(w, "Invalid output format", http.StatusInternalServerError)
		return
	}
	responseCode := 200
	if status, ok := outputMap["status"].(float64); ok {
		responseCode = int(status)
	}
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helpers.FillResponseHeaders(responseHeaders, w)

	w.WriteHeader(responseCode)

	w.Write(responseBody)
}

func main() {
	helpers.SetupLogging()

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

	logrus.Fatal(http.ListenAndServe(":8080", nil))
}
