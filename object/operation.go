package object

import (
	"context"
	"fmt"

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
	output, ok := stepMap["output"].(map[string]interface{})
	if !ok {
		err := fmt.Errorf("invalid output format")
		return err.Error(), "error", err
	}

	logrus.Infof("output: %v", output)
	logrus.Infof("next: %v", next)

	body, err := helpers.TransformBody(stepOutputs, output)

	if err != nil {
		logrus.Errorf("could not transform body: %v", err)
		return err.Error(), "error", err
	}

	return body, next, nil
}
