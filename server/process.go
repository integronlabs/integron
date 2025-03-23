package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/integronlabs/integron/helpers"

	"github.com/sirupsen/logrus"
)

func (s *Server) ProcessStep(currentStepKey string, w http.ResponseWriter, steps map[string]interface{}, stepOutputs map[string]interface{}, stepInput interface{}) (interface{}, string) {
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

	stepOutput, next, err := handler(s.Ctx, stepMap, stepOutputs)
	if err != nil {
		if stepType == "error" {
			message, _ := json.Marshal(map[string]interface{}{"message": stepInput})
			http.Error(w, string(message), http.StatusInternalServerError)
			return nil, "end"
		}
		return err.Error(), "error"
	}
	logrus.Debugf("Step %s completed", currentStepKey)
	logrus.Debugf("Step outputs: %v", stepOutput)
	return stepOutput, next
}
