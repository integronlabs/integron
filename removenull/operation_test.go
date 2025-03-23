package removenull

import (
	"context"
	"testing"
)

const EXPECTED_NIL_GOT = "Expected nil, got %v"
const EXPECTED_ERROR_GOT_NIL = "Expected error, got nil"
const EXPECTED_BUT_GOT = "Expected %v, got %v"
const VALID_INPUT = "$.output"

var validOutputMap = map[string]interface{}{
	"output": interface{}{
		map[string]interface{}{
			"message": "world",
		},
	},
}

func TestRun(t *testing.T) {
	ctx := context.Background()
	stepMap := map[string]interface{}{
		"next":  "next",
		"input": VALID_INPUT,
	}
	stepOutputs := validOutputMap

	expectedOutput := interface{}{
		map[string]interface{}{
			"message": "world",
		},
	}

	expectedNext := "next"

	output, next, err := Run(ctx, stepMap, stepOutputs)

	if err != nil {
		t.Errorf(EXPECTED_NIL_GOT, err)
	}
	if output.([]interface{})[0].(map[string]interface{})["message"] != expectedOutput[0].(map[string]interface{})["message"] {
		t.Errorf(EXPECTED_BUT_GOT, expectedOutput, output)
	}
	if next != expectedNext {
		t.Errorf(EXPECTED_BUT_GOT, expectedNext, next)
	}
}

func TestRunInvalidNext(t *testing.T) {
	ctx := context.Background()
	stepMap := map[string]interface{}{
		"next":  1,
		"input": VALID_INPUT,
	}
	stepOutputs := validOutputMap

	expectedError := "invalid next format"

	_, _, err := Run(ctx, stepMap, stepOutputs)

	if err == nil {
		t.Error(EXPECTED_ERROR_GOT_NIL)
	}
	if err.Error() != expectedError {
		t.Errorf(EXPECTED_BUT_GOT, expectedError, err)
	}
}

func TestRunInvalidInput(t *testing.T) {
	ctx := context.Background()
	stepMap := map[string]interface{}{
		"next":  "next",
		"input": 1,
	}
	stepOutputs := validOutputMap

	expectedError := "invalid input format"

	_, _, err := Run(ctx, stepMap, stepOutputs)

	if err == nil {
		t.Error(EXPECTED_ERROR_GOT_NIL)
	}
	if err.Error() != expectedError {
		t.Errorf(EXPECTED_BUT_GOT, expectedError, err)
	}
}

func TestRunInvalidInputPath(t *testing.T) {
	ctx := context.Background()
	stepMap := map[string]interface{}{
		"next":  "next",
		"input": "$.invalid",
	}
	stepOutputs := validOutputMap

	expectedError := "unknown key invalid"

	_, _, err := Run(ctx, stepMap, stepOutputs)

	if err == nil {
		t.Error(EXPECTED_ERROR_GOT_NIL)
	}
	if err.Error() != expectedError {
		t.Errorf(EXPECTED_BUT_GOT, expectedError, err)
	}
}
