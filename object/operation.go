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

	logrus.WithContext(ctx).Debugf("output: %v", output)
	logrus.WithContext(ctx).Debugf("next: %v", next)

	body := helpers.TransformBody(stepOutputs, output)

	return body, next, nil
}
