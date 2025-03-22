package helpers

import (
	"fmt"
	"regexp"
	"strings"

	"log"

	"github.com/PaesslerAG/jsonpath"
)

func Replace(input string, stepOutputs map[string]interface{}) (string, error) {
	re := regexp.MustCompile(`\$\.[a-zA-Z0-9_\[\]\.]+`)
	matches := re.FindAllString(input, -1)

	for _, match := range matches {
		value, err := jsonpath.Get(match, stepOutputs)
		if err != nil {
			return input, err
		}
		input = strings.ReplaceAll(input, match, fmt.Sprintf("%v", value))
	}

	return input, nil
}

func TransformBody(body interface{}, output interface{}) interface{} {

	// if output is array, go through each element and transform
	if outputArray, ok := output.([]interface{}); ok {

		transformedBody := make([]interface{}, 0)
		for _, outputMap := range outputArray {
			transformed := TransformBody(body, outputMap)
			transformedBody = append(transformedBody, transformed)
		}
		return transformedBody
	}

	if outputMap, ok := output.(map[string]interface{}); ok {
		transformedBody := make(map[string]interface{})
		// if output is not array, transform
		for key, value := range outputMap {
			transformedBody[key] = TransformBody(body, value)
		}
		return transformedBody
	}

	if valueStr, ok := output.(string); ok {
		if strings.HasPrefix(valueStr, "$") {
			// get value from body
			value, err := jsonpath.Get(valueStr, body)
			if err != nil {
				log.Printf("could not read value from body: %v", err)
			}
			return value
		} else {
			return valueStr
		}
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
