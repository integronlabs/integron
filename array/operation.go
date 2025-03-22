package http

import (
	"fmt"

	"context"

	"github.com/PaesslerAG/jsonpath"
	"github.com/integronlabs/integron/helpers"
	"github.com/sirupsen/logrus"
)

func Run(ctx context.Context, stepMap map[string]interface{}, stepOutputs map[string]interface{}) (interface{}, string, error) {
	// get values

	next := stepMap["next"].(string)
	inputString := stepMap["input"].(string)
	output := stepMap["output"].(map[string]interface{})

	logrus.Infof("inputString: %v", inputString)
	logrus.Infof("output: %v", output)
	logrus.Infof("next: %v", next)

	// replace placeholders in input
	inputMap, err := jsonpath.Get(inputString, stepOutputs)
	if err != nil {
		logrus.Errorf("could not read value from input: %v", err)
		return err.Error(), next, err
	}

	logrus.Infof("inputMap: %v", inputMap)

	inputArray, ok := inputMap.([]interface{})
	if !ok {
		err := fmt.Errorf("invalid input format")
		return err.Error(), "error", err
	}

	body, err := helpers.TransformArray(inputArray, output)

	if err != nil {
		logrus.Errorf("could not transform body: %v", err)
		return body, next, err
	}

	return body, next, nil
}
