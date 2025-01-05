package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	httpOperation "github.com/integronlabs/integron/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/swaggest/swgui/v5emb"
)

var router routers.Router
var ctx context.Context

func handler(w http.ResponseWriter, r *http.Request) {
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
	stepOutputs := make(map[string]interface{})
	input := make(map[string]interface{})
	// path params
	for key, value := range pathParams {
		input[key] = value
	}
	// query params
	for key, value := range r.URL.Query() {
		input[key] = value[0]
	}

	stepOutputs["request"] = input

	steps, ok := route.PathItem.GetOperation(route.Method).Extensions["x-integron-steps"].([]interface{})

	if !ok {
		http.Error(w, "Invalid x-integron-steps", http.StatusInternalServerError)
		return
	}
	for _, step := range steps {
		stepMap, ok := step.(map[string]interface{})
		if !ok {
			http.Error(w, "Invalid step definition", http.StatusInternalServerError)
			return
		}

		switch (stepMap["type"]).(string) {
		case "http":
			stepOutputs[stepMap["name"].(string)], err = httpOperation.Run(stepMap, input, stepOutputs)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		output = stepOutputs[stepMap["name"].(string)]
	}
	responseHeaders := http.Header{"Content-Type": []string{"application/json"},
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

	for k, v := range responseHeaders {
		w.Header().Set(k, v[0])
	}

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

	http.ListenAndServe(":8080", nil)
}
