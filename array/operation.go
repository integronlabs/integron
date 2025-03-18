package http

import (
	"log"

	"context"

	"github.com/PaesslerAG/jsonpath"
	"github.com/integronlabs/integron/helpers"
)

func Run(ctx context.Context, stepMap map[string]interface{}, input map[string]interface{}, stepOutputs map[string]interface{}) (interface{}, string, error) {
	// get values

	next := stepMap["next"].(string)
	inputString := stepMap["input"].(string)
	output := stepMap["output"].(map[string]interface{})

	log.Printf("inputString: %v", inputString)
	log.Printf("output: %v", output)
	log.Printf("next: %v", next)

	// replace placeholders in input
	inputMap, err := jsonpath.Get(inputString, stepOutputs)
	if err != nil {
		log.Printf("could not read value from input: %v", err)
		return err.Error(), next, err
	}

	log.Printf("inputMap: %v", inputMap)

	body := helpers.TransformBody(inputMap, output)

	return body, next, nil
}
