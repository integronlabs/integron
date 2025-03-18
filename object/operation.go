package http

import (
	"log"

	"context"

	"github.com/integronlabs/integron/helpers"
)

func Run(ctx context.Context, stepMap map[string]interface{}, input map[string]interface{}, stepOutputs map[string]interface{}) (interface{}, string, error) {
	// get values

	next := stepMap["next"].(string)
	output := stepMap["output"].(map[string]interface{})

	log.Printf("output: %v", output)
	log.Printf("next: %v", next)

	body := helpers.TransformBody(stepOutputs, output)

	return body, next, nil
}
