package http

import (
	"log"
	"strings"

	"context"

	"github.com/PaesslerAG/jsonpath"
)

func transformBody(body interface{}, output map[string]interface{}) interface{} {
	// if output is array, go through each element and transform
	if bodyArray, ok := body.([]interface{}); ok {
		log.Println("body is array")
		transformedBody := make([]interface{}, 0)
		for _, bodyMap := range bodyArray {
			log.Println("bodyMap: ", bodyMap)
			transformed := transformBody(bodyMap, output)
			log.Printf("transformed: %v", transformed)
			transformedBody = append(transformedBody, transformed)
		}
		return transformedBody
	}

	transformedBody := make(map[string]interface{})
	// if output is not array, transform
	for key, value := range output {
		if valueStr, ok := value.(string); ok {
			if strings.HasPrefix(valueStr, "$.") {
				// get value from body
				value, err := jsonpath.Get(valueStr, body)
				if err != nil {
					log.Printf("could not read value from body: %v", err)
				}
				transformedBody[key] = value
			} else {
				transformedBody[key] = value
			}
		}
	}
	return transformedBody
}

func Run(ctx context.Context, stepMap map[string]interface{}, input map[string]interface{}, stepOutputs map[string]interface{}) (interface{}, string, error) {
	// get values

	next := stepMap["next"].(string)
	output := stepMap["output"].(map[string]interface{})

	log.Printf("output: %v", output)
	log.Printf("next: %v", next)

	body := transformBody(stepOutputs, output)

	return body, next, nil
}
