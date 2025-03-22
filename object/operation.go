package object

import (
	"context"

	"github.com/integronlabs/integron/helpers"
	"github.com/sirupsen/logrus"
)

func Run(ctx context.Context, stepMap map[string]interface{}, stepOutputs map[string]interface{}) (interface{}, string, error) {
	// get values

	next := stepMap["next"].(string)
	output := stepMap["output"].(map[string]interface{})

	logrus.Infof("output: %v", output)
	logrus.Infof("next: %v", next)

	body, err := helpers.TransformBody(stepOutputs, output)

	if err != nil {
		logrus.Errorf("could not transform body: %v", err)
		return body, next, err
	}

	return body, next, nil
}
