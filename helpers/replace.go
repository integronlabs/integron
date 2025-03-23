package helpers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PaesslerAG/jsonpath"
)

func Replace(input string, stepOutputs interface{}) string {
	re := regexp.MustCompile(`\$\.[a-zA-Z0-9_\[\]\.]+`)
	matches := re.FindAllString(input, -1)

	for _, match := range matches {
		value, _ := jsonpath.Get(match, stepOutputs)
		input = strings.ReplaceAll(input, match, fmt.Sprintf("%v", value))
	}

	return input
}

func transformBodyArray(outputArray []interface{}, body interface{}) []interface{} {
	transformedBody := make([]interface{}, 0)
	for _, outputMap := range outputArray {
		transformed := TransformBody(body, outputMap)
		transformedBody = append(transformedBody, transformed)
	}
	return transformedBody
}

func transformBodyMap(outputMap map[string]interface{}, body interface{}) map[string]interface{} {
	transformedBody := make(map[string]interface{})
	for key, value := range outputMap {
		transformed := TransformBody(body, value)
		transformedBody[key] = transformed
	}
	return transformedBody
}

func transformBodyString(output string, body interface{}) interface{} {
	if strings.HasPrefix(output, "$") {
		// get value from body
		value, _ := jsonpath.Get(output, body)
		return value
	} else {
		value := Replace(output, body)
		return value
	}
}

func TransformBody(body interface{}, output interface{}) interface{} {

	// if output is array, go through each element and transform
	if outputArray, ok := output.([]interface{}); ok {
		return transformBodyArray(outputArray, body)
	}

	if outputMap, ok := output.(map[string]interface{}); ok {
		return transformBodyMap(outputMap, body)
	}

	if valueStr, ok := output.(string); ok {
		return transformBodyString(valueStr, body)
	}
	return output
}

func TransformArray(inputArray []interface{}, output map[string]interface{}) []interface{} {
	transformedArray := make([]interface{}, 0)
	for _, inputMap := range inputArray {
		transformed := TransformBody(inputMap, output)
		transformedArray = append(transformedArray, transformed)
	}
	return transformedArray
}
