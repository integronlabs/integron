package server

import (
	"fmt"
	"net/http"

	"github.com/integronlabs/integron/helpers"

	"github.com/sirupsen/logrus"
)

func (s *Server) ProcessStep(r *http.Request, currentStepKey string, w http.ResponseWriter, steps map[string]interface{}, stepOutputs map[string]interface{}, stepInput interface{}) (interface{}, string) {
	ctx := r.Context()

	logrus.WithContext(ctx).Debugf("Processing step: %s", currentStepKey)

	var next string
	var err error

	step, ok := steps[currentStepKey]
	if !ok {
		return fmt.Errorf(helpers.INVALID_STEP_DEFINITION), "error"
	}
	stepMap, ok := step.(map[string]interface{})
	if !ok {
		return fmt.Errorf(helpers.INVALID_STEP_DEFINITION), "error"
	}

	stepType, ok := stepMap["type"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid step type"), "error"
	}

	handler, err := GetStepHandler(stepType)
	if err != nil {
		return fmt.Errorf("unknown step type: %s", stepType), "error"
	}

	if stepType == "error" {
		Error(r, w, stepInput.(string), http.StatusInternalServerError, "EXCEPTION")
		return nil, "end"
	}

	stepOutput, next, err := handler(ctx, stepMap, stepOutputs)
	if err != nil {
		return err.Error(), "error"
	}
	logrus.WithContext(ctx).Debugf("Step %s completed", currentStepKey)
	logrus.WithContext(ctx).Debugf("Step outputs: %v", stepOutput)
	return stepOutput, next
}
