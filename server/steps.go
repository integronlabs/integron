package server

import "errors"

var stepRegistry = make(map[string]StepHandler)

// RegisterStep registers a step handler for a specific type.
func RegisterStep(stepType string, handler StepHandler) {
	stepRegistry[stepType] = handler
}

// GetStepHandler retrieves a step handler by type.
func GetStepHandler(stepType string) (StepHandler, error) {
	handler, exists := stepRegistry[stepType]
	if !exists {
		return nil, errors.New("step type not registered")
	}
	return handler, nil
}
