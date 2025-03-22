package helpers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PaesslerAG/jsonpath"
)

func Replace(input string, stepOutputs interface{}) (string, error) {
	re := regexp.MustCompile(`\$\.[a-zA-Z0-9_\[\]\.]+`)
	matches := re.FindAllString(input, -1)

	for _, match := range matches {
		value, err := jsonpath.Get(match, stepOutputs)
		if err != nil {
			return err.Error(), err
		}
		input = strings.ReplaceAll(input, match, fmt.Sprintf("%v", value))
	}

	return input, nil
}

func TransformBody(body interface{}, output interface{}) (interface{}, error) {

	// if output is array, go through each element and transform
	if outputArray, ok := output.([]interface{}); ok {

		transformedBody := make([]interface{}, 0)
		for _, outputMap := range outputArray {
			transformed, err := TransformBody(body, outputMap)
			if err != nil {
				return transformedBody, err
			}
			transformedBody = append(transformedBody, transformed)
		}
		return transformedBody, nil
	}

	if outputMap, ok := output.(map[string]interface{}); ok {
		transformedBody := make(map[string]interface{})
		// if output is not array, transform
		for key, value := range outputMap {
			transformed, err := TransformBody(body, value)
			if err != nil {
				return transformedBody, err
			}
			transformedBody[key] = transformed
		}
		return transformedBody, nil
	}

	if valueStr, ok := output.(string); ok {
		if strings.HasPrefix(valueStr, "$") {
			// get value from body
			value, err := jsonpath.Get(valueStr, body)
			return value, err
		} else {
			value, err := Replace(valueStr, body)
			return value, err
		}
	}
	return output, nil
}

func TransformArray(inputArray []interface{}, output map[string]interface{}) ([]interface{}, error) {
	transformedArray := make([]interface{}, 0)
	for _, inputMap := range inputArray {
		transformed, err := TransformBody(inputMap, output)
		if err != nil {
			return transformedArray, err
		}
		transformedArray = append(transformedArray, transformed)
	}
	return transformedArray, nil
}
