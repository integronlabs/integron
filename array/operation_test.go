package array

import (
	"context"
	"testing"
)

const EXPECTED_NIL_GOT = "Expected nil, got %v"
const EXPECTED_ERROR_GOT_NIL = "Expected error, got nil"
const EXPECTED_BUT_GOT = "Expected %v, got %v"
const VALID_OUTPUT = "$.message"
const VALID_INPUT = "$.output"

var validOutputMap = map[string]interface{}{
	"output": []interface{}{
		map[string]interface{}{
			"message": "world",
		},
	},
}

func TestRun(t *testing.T) {
	ctx := context.Background()
	stepMap := map[string]interface{}{
		"next": "next",
		"output": map[string]interface{}{
			"message": VALID_OUTPUT,
		},
		"input": VALID_INPUT,
	}
	stepOutputs := validOutputMap

	expectedOutput := []interface{}{
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

func TestRunInvalidNextFormat(t *testing.T) {
	ctx := context.Background()
	stepMap := map[string]interface{}{
		"next": 1,
		"output": map[string]interface{}{
			"message": VALID_OUTPUT,
		},
		"input": VALID_INPUT,
	}
	stepOutputs := validOutputMap

	expectedError := "invalid next format"

	output, next, err := Run(ctx, stepMap, stepOutputs)

	if err == nil {
		t.Error(EXPECTED_ERROR_GOT_NIL)
	}
	if output != expectedError {
		t.Errorf(EXPECTED_BUT_GOT, expectedError, output)
	}
	if next != "error" {
		t.Errorf(EXPECTED_BUT_GOT, "error", next)
	}
}

func TestRunInvalidInputFormat(t *testing.T) {
	ctx := context.Background()
	stepMap := map[string]interface{}{
		"next": "next",
		"output": map[string]interface{}{
			"message": VALID_OUTPUT,
		},
		"input": 1,
	}
	stepOutputs := validOutputMap

	expectedError := "invalid input format"

	output, next, err := Run(ctx, stepMap, stepOutputs)

	if err == nil {
		t.Error(EXPECTED_ERROR_GOT_NIL)
	}
	if output != expectedError {
		t.Errorf(EXPECTED_BUT_GOT, expectedError, output)
	}
	if next != "error" {
		t.Errorf(EXPECTED_BUT_GOT, "error", next)
	}
}

func TestRunInvalidOutputFormat(t *testing.T) {
	ctx := context.Background()
	stepMap := map[string]interface{}{
		"next":   "next",
		"output": 1,
		"input":  VALID_INPUT,
	}
	stepOutputs := validOutputMap

	expectedError := "invalid output format"

	output, next, err := Run(ctx, stepMap, stepOutputs)

	if err == nil {
		t.Error(EXPECTED_ERROR_GOT_NIL)
	}
	if output != expectedError {
		t.Errorf(EXPECTED_BUT_GOT, expectedError, output)
	}
	if next != "error" {
		t.Errorf(EXPECTED_BUT_GOT, "error", next)
	}
}

func TestRunInvalidInputFormat2(t *testing.T) {
	ctx := context.Background()
	stepMap := map[string]interface{}{
		"next": "next",
		"output": map[string]interface{}{
			"message": VALID_OUTPUT,
		},
		"input": VALID_INPUT,
	}
	stepOutputs := map[string]interface{}{
		"output": 1,
	}

	expectedError := "invalid input format"

	output, next, err := Run(ctx, stepMap, stepOutputs)

	if err == nil {
		t.Error(EXPECTED_ERROR_GOT_NIL)
	}
	if output != expectedError {
		t.Errorf(EXPECTED_BUT_GOT, expectedError, output)
	}
	if next != "error" {
		t.Errorf(EXPECTED_BUT_GOT, "error", next)
	}
}

func TestRunInvalidInputFormat3(t *testing.T) {
	ctx := context.Background()
	stepMap := map[string]interface{}{
		"next": "next",
		"output": map[string]interface{}{
			"message": VALID_OUTPUT,
		},
		"input": "$.output[]",
	}
	stepOutputs := map[string]interface{}{
		"output": 1,
	}

	expectedError := "parsing error: $.output[]	:1:10 - 1:11 unexpected \"]\" while scanning extensions"

	output, next, err := Run(ctx, stepMap, stepOutputs)

	if err == nil {
		t.Error(EXPECTED_ERROR_GOT_NIL)
	}
	if output != expectedError {
		t.Errorf(EXPECTED_BUT_GOT, expectedError, output)
	}
	if next != "error" {
		t.Errorf(EXPECTED_BUT_GOT, "error", next)
	}
}
