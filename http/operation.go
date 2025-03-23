package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/integronlabs/integron/helpers"
)

func getActions(responsesMap map[string]interface{}, statusCodeStr string) (map[string]interface{}, string, error) {
	statusMap, ok := responsesMap[statusCodeStr].(map[string]interface{})
	if !ok {
		statusMap, ok = responsesMap["default"].(map[string]interface{})
		if !ok {
			return nil, "error", fmt.Errorf("could not find actions for status %s", statusCodeStr)
		}
	}
	outputMap, ok := statusMap["output"].(map[string]interface{})
	if !ok {
		return nil, "error", fmt.Errorf("invalid output format")
	}
	next, ok := statusMap["next"].(string)
	if !ok {
		return nil, "error", fmt.Errorf("invalid next format")
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
	response, err := client.Do(httpRequest)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func Run(ctx context.Context, client *http.Client, stepMap map[string]interface{}, stepOutputs map[string]interface{}) (interface{}, string, error) {
	// get values
	method, _ := stepMap["method"].(string)
	url, _ := stepMap["url"].(string)
	requestBodyMap, _ := stepMap["body"].(map[string]interface{})
	headers, _ := stepMap["headers"].(map[string]interface{})
	responsesMap, _ := stepMap["responses"].(map[string]interface{})

	requestBody := helpers.TransformBody(stepOutputs, requestBodyMap)

	requestBodyJson, _ := json.Marshal(requestBody)
	requestBodyString := string(requestBodyJson)

	response, err := httpRequest(ctx, client, method, url, requestBodyString, headers, stepOutputs)

	if err != nil {
		return err.Error(), "error", err
	}

	defer response.Body.Close()

	var responseData interface{}
	if err := json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		return err.Error(), "error", err
	}

	responseMap := map[string]interface{}{
		"status":  response.StatusCode,
		"headers": response.Header,
		"body":    responseData,
	}

	statusCodeStr := fmt.Sprintf("%d", response.StatusCode)

	outputMap, next, err := getActions(responsesMap, statusCodeStr)

	if err != nil {
		return err.Error(), "error", err
	}

	body := helpers.TransformBody(responseMap, outputMap)

	return body, next, nil
}
