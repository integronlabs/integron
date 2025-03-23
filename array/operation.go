package array

import (
	"fmt"

	"context"

	"github.com/PaesslerAG/jsonpath"
	"github.com/integronlabs/integron/helpers"
	"github.com/sirupsen/logrus"
)

func Run(ctx context.Context, stepMap map[string]interface{}, stepOutputs map[string]interface{}) (interface{}, string, error) {
	// get values

	next, ok := stepMap["next"].(string)
	if !ok {
		err := fmt.Errorf("invalid next format")
		return err.Error(), "error", err
	}
	inputString, ok := stepMap["input"].(string)
	if !ok {
		err := fmt.Errorf("invalid input format")
		return err.Error(), "error", err
	}
	output, ok := stepMap["output"].(map[string]interface{})
	if !ok {
		err := fmt.Errorf("invalid output format")
		return err.Error(), "error", err
	}

	logrus.Debugf("inputString: %v", inputString)
	logrus.Debugf("output: %v", output)
	logrus.Debugf("next: %v", next)

	// replace placeholders in input
	inputMap, err := jsonpath.Get(inputString, stepOutputs)
	if err != nil {
		logrus.Errorf("could not read value from input: %v", err)
		return err.Error(), "error", err
	}

	logrus.Debugf("inputMap: %v", inputMap)

	inputArray, ok := inputMap.([]interface{})
	if !ok {
		err := fmt.Errorf("invalid input format")
		return err.Error(), "error", err
	}

	body, err := helpers.TransformArray(inputArray, output)

	if err != nil {
		logrus.Errorf("could not transform body: %v", err)
		return err.Error(), "error", err
	}

	return body, next, nil
}
