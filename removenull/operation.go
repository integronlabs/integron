package removenull

import (
	"fmt"

	"context"

	"github.com/PaesslerAG/jsonpath"
	"github.com/integronlabs/integron/helpers"
	"github.com/sirupsen/logrus"
)

func Run(ctx context.Context, stepMap map[string]interface{}, stepOutputs map[string]interface{}) (interface{}, string, error) {
	// get values

	inputString, ok := stepMap["input"].(string)
	if !ok {
		err := fmt.Errorf("invalid input format")
		return err.Error(), "error", err
	}
	next, ok := stepMap["next"].(string)
	if !ok {
		err := fmt.Errorf("invalid next format")
		return err.Error(), "error", err
	}

	logrus.WithContext(ctx).Debugf("inputString: %v", inputString)
	logrus.WithContext(ctx).Debugf("next: %v", next)

	// replace placeholders in input
	inputMap, err := jsonpath.Get(inputString, stepOutputs)
	if err != nil {
		logrus.Errorf("could not read value from input: %v", err)
		return err.Error(), "error", err
	}

	logrus.WithContext(ctx).Debugf("inputMap: %v", inputMap)

	body := helpers.RemoveNull(inputMap)

	return body, next, nil
}
