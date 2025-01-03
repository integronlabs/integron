package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/swaggest/swgui/v5emb"
)

var router routers.Router
var ctx context.Context

func init() {
	ctx = context.Background()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	doc, err := loader.LoadFromFile("../docs/openapi.yaml")
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
}

type Output struct {
	Body map[string]interface{} `json:"-"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
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

	_ = openapi3filter.ValidateRequest(ctx, requestValidationInput)

	// Handle that request
	var steps []interface{}
	var ok bool
	var output Output
	input := make(map[string]interface{})
	// path params
	for key, value := range pathParams {
		input[key] = value
	}
	// query params
	for key, value := range r.URL.Query() {
		input[key] = value[0]
	}
	if route.Method == http.MethodGet {
		steps, ok = route.PathItem.Get.Extensions["x-integron-steps"].([]interface{})
	}
	if route.Method == http.MethodPost {
		steps, ok = route.PathItem.Post.Extensions["x-integron-steps"].([]interface{})
	}
	if route.Method == http.MethodPut {
		steps, ok = route.PathItem.Put.Extensions["x-integron-steps"].([]interface{})
	}
	if route.Method == http.MethodDelete {
		steps, ok = route.PathItem.Delete.Extensions["x-integron-steps"].([]interface{})
	}
	if !ok {
		http.Error(w, "Invalid x-integron-steps format", http.StatusInternalServerError)
		return
	}
	for _, step := range steps {
		stepMap, ok := step.(map[string]interface{})
		if !ok {
			http.Error(w, "Invalid step format", http.StatusInternalServerError)
			return
		}

		if stepMap["type"] == "http" {
			var response *http.Response
			// call url with method
			method, _ := stepMap["method"].(string)
			url, _ := stepMap["url"].(string)
			responseType, _ := stepMap["response"].(string)
			if method == http.MethodGet {
				// replace input in url
				for key, value := range input {
					url = strings.ReplaceAll(url, "input."+key, value.(string))
				}
				for key, value := range output.Body {
					url = strings.ReplaceAll(url, "output."+key, value.(string))
				}
				fmt.Println(url, input, output.Body)
				response, err = http.Get(url)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
			defer response.Body.Close()
			// json decode response
			if responseType == "array" {
				var responseData []map[string]interface{}
				if err := json.NewDecoder(response.Body).Decode(&responseData); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				if len(responseData) == 0 {
					http.Error(w, "Empty response", http.StatusInternalServerError)
					return
				}

				outputMap, ok := stepMap["output"].(map[string]interface{})
				if !ok {
					http.Error(w, "Invalid output format", http.StatusInternalServerError)
					return
				}

				output.Body = make(map[string]interface{})
				for key, path := range outputMap {
					if value, found := getValueFromPath(responseData[0], path.(string)); found {
						output.Body[key] = value
					}

				}
			} else {
				var responseData map[string]interface{}
				if err := json.NewDecoder(response.Body).Decode(&responseData); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				outputMap, ok := stepMap["output"].(map[string]interface{})
				if !ok {
					http.Error(w, "Invalid output format", http.StatusInternalServerError)
					return
				}

				output.Body = make(map[string]interface{})
				for key, path := range outputMap {
					if value, found := getValueFromPath(responseData, path.(string)); found {
						output.Body[key] = value
					}

				}
			}
		}
	}
	responseHeaders := http.Header{"Content-Type": []string{"application/json"}}
	responseCode := 200
	jsonBody, _ := json.Marshal(output.Body)
	fmt.Println(string(jsonBody))
	responseBody := []byte(jsonBody)

	// Validate response
	responseValidationInput := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: requestValidationInput,
		Status:                 responseCode,
		Header:                 responseHeaders,
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

	// CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	w.WriteHeader(responseCode)

	w.Write(responseBody)
}

func getValueFromPath(data map[string]interface{}, path string) (interface{}, bool) {
	keys := strings.Split(path, ".")
	var value interface{} = data
	for _, key := range keys {
		if m, ok := value.(map[string]interface{}); ok {
			value = m[key]
		} else {
			return nil, false
		}
	}
	return value, true
}

func main() {
	http.Handle("/api/", http.StripPrefix("/api", http.HandlerFunc(handler)))

	fs := http.FileServer(http.Dir("../docs/"))
	http.Handle("/docs/", http.StripPrefix("/docs/", fs))

	http.Handle("/", v5emb.New(
		"Integron Sunrise",
		"/docs/openapi.yaml",
		"/",
	))

	http.ListenAndServe(":8080", nil)
}
