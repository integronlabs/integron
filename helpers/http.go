package helpers

import (
	"net/http"
)

func ExtractParams(pathParams map[string]string, queryParams map[string][]string) map[string]interface{} {
	params := make(map[string]interface{})
	for key, value := range pathParams {
		params[key] = value
	}
	for key, value := range queryParams {
		params[key] = value[0]
	}
	return params
}

func FillResponseHeaders(responseHeaders http.Header, w http.ResponseWriter) {
	for k, v := range responseHeaders {
		w.Header().Set(k, v[0])
	}
}
