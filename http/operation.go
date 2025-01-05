package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/yalp/jsonpath"

	"github.com/integronlabs/integron/helpers"
)

func Run(stepMap map[string]interface{}, input map[string]interface{}, stepOutputs map[string]interface{}) (interface{}, error) {
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
		return body, err
	}
	requestBodyString, err = helpers.Replace(requestBodyString, stepOutputs)
	if err != nil {
		return body, err
	}

	fmt.Println(url, requestBodyString)

	client := &http.Client{}

	httpRequest, err := http.NewRequest(method, url, strings.NewReader(requestBodyString))
	if err != nil {
		return body, err
	}
	// set headers
	for key, value := range headers {
		value, err := helpers.Replace(value.(string), stepOutputs)
		if err != nil {
			return body, err
		}
		httpRequest.Header.Set(key, value)
	}
	response, err := client.Do(httpRequest)

	if err != nil {
		return body, err
	}

	defer response.Body.Close()

	var responseData interface{}
	if err := json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		return body, err
	}

	outputMap, ok := stepMap["output"].(map[string]interface{})
	if !ok {
		return body, fmt.Errorf("invalid output format")
	}

	for key, path := range outputMap {
		if value, err := jsonpath.Read(responseData, path.(string)); err == nil {
			body[key] = value
		}
	}

	return body, nil
}
