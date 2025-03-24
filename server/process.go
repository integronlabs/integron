package server

import (
	"fmt"
	"net/http"

	"github.com/integronlabs/integron/helpers"

	"github.com/sirupsen/logrus"
)

func (s *Server) ProcessStep(r *http.Request, currentStepKey string, w http.ResponseWriter, steps map[string]interface{}, stepOutputs map[string]interface{}, stepInput interface{}) (interface{}, string) {
	ctx := r.Context()

	logrus.Debugf("Processing step: %s", currentStepKey)

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

	stepOutput, next, err := handler(ctx, stepMap, stepOutputs)
	if err != nil {
		if stepType == "error" {
			Error(r, w, stepInput.(error).Error(), http.StatusInternalServerError, "EXCEPTION")
			return nil, "end"
		}
		return err.Error(), "error"
	}
	logrus.Debugf("Step %s completed", currentStepKey)
	logrus.Debugf("Step outputs: %v", stepOutput)
	return stepOutput, next
}
