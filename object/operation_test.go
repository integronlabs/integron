package object

import (
	"context"
	"testing"
)

const EXPECTED_NIL_GOT = "Expected nil, got %v"
const EXPECTED_ERROR_GOT_NIL = "Expected error, got nil"
const EXPECTED_BUT_GOT = "Expected %v, got %v"
const VALID_OUTPUT = "$.output.message"

var validOutputMap = map[string]interface{}{
	"output": map[string]interface{}{
		"message": "world",
	},
}

func TestRun(t *testing.T) {
	ctx := context.Background()
	stepMap := map[string]interface{}{
		"next": "next",
		"output": map[string]interface{}{
			"message": VALID_OUTPUT,
		},
	}
	stepOutputs := validOutputMap

	expectedOutput := map[string]interface{}{
		"message": "world",
	}

	expectedNext := "next"

	output, next, err := Run(ctx, stepMap, stepOutputs)

	if err != nil {
		t.Errorf(EXPECTED_NIL_GOT, err)
	}
	if output.(map[string]interface{})["message"] != expectedOutput["message"] {
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
	}
	stepOutputs := validOutputMap
	expectedOutput := "invalid next format"

	output, next, err := Run(ctx, stepMap, stepOutputs)

	if err == nil {
		t.Error(EXPECTED_ERROR_GOT_NIL)
	}
	if output != expectedOutput {
		t.Errorf(EXPECTED_BUT_GOT, expectedOutput, output)
	}
	if next != "error" {
		t.Errorf(EXPECTED_BUT_GOT, "error", next)
	}
}

func TestRunInvalidOutputFormat(t *testing.T) {
	ctx := context.Background()
	stepMap := map[string]interface{}{
		"next":   "next",
		"output": "invalid",
	}
	stepOutputs := validOutputMap
	expectedOutput := "invalid output format"
	output, next, err := Run(ctx, stepMap, stepOutputs)

	if err == nil {
		t.Error(EXPECTED_ERROR_GOT_NIL)
	}
	if output != expectedOutput {
		t.Errorf(EXPECTED_BUT_GOT, expectedOutput, output)
	}
	if next != "error" {
		t.Errorf(EXPECTED_BUT_GOT, "error", next)
	}
}

func TestRunCouldNotTransformBody(t *testing.T) {
	ctx := context.Background()
	stepMap := map[string]interface{}{
		"next": "next",
		"output": map[string]interface{}{
			"message": "$.output[]",
		},
	}
	stepOutputs := map[string]interface{}{
		"output": "invalid",
	}

	expectedOutput := "parsing error: $.output[]	:1:10 - 1:11 unexpected \"]\" while scanning extensions"

	output, next, err := Run(ctx, stepMap, stepOutputs)

	if err == nil {
		t.Error(EXPECTED_ERROR_GOT_NIL)
	}
	if output != expectedOutput {
		t.Errorf(EXPECTED_BUT_GOT, expectedOutput, output)
	}
	if next != "error" {
		t.Errorf(EXPECTED_BUT_GOT, "error", next)
	}
}
