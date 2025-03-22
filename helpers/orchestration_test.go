package helpers

import (
	"testing"
)

const EXPECTED_ERROR_GOT_NIL = "Expected error, got nil"
const EXPECTED_NIL_GOT = "Expected nil, got %v"

func TestCreateStepsMap(t *testing.T) {
	stepsArray := []interface{}{
		map[string]interface{}{
			"name": "step1",
		},
		map[string]interface{}{
			"name": "step2",
		},
	}
	steps, err := CreateStepsMap(stepsArray)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(steps) != 2 {
		t.Errorf("Expected 2 steps, got %d", len(steps))
	}
}

func TestCreateStepsMapInvalidStepDefinition(t *testing.T) {
	stepsArray := []interface{}{
		map[string]interface{}{
			"name": "step1",
		},
		"invalid",
	}
	steps, err := CreateStepsMap(stepsArray)
	if err == nil {
		t.Error(EXPECTED_ERROR_GOT_NIL)
	}
	if steps != nil {
		t.Errorf(EXPECTED_NIL_GOT, steps)
	}
}

func TestCreateStepsMapInvalidStepDefinition2(t *testing.T) {
	stepsArray := []interface{}{
		"invalid",
	}
	steps, err := CreateStepsMap(stepsArray)
	if err == nil {
		t.Error(EXPECTED_ERROR_GOT_NIL)
	}
	if steps != nil {
		t.Errorf(EXPECTED_NIL_GOT, steps)
	}
}

func TestCreateStepsMapInvalidStepDefinition3(t *testing.T) {
	stepsArray := []interface{}{}
	steps, err := CreateStepsMap(stepsArray)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(steps) != 0 {
		t.Errorf("Expected 0 steps, got %d", len(steps))
	}
}

func TestCreateStepsMapInvalidStepDefinition4(t *testing.T) {
	stepsArray := []interface{}{
		"invalid",
		map[string]interface{}{
			"name": "step1",
		},
	}
	steps, err := CreateStepsMap(stepsArray)
	if err == nil {
		t.Error(EXPECTED_ERROR_GOT_NIL)
	}
	if steps != nil {
		t.Errorf(EXPECTED_NIL_GOT, steps)
	}
}

func TestCreateStepsMapInvalidStepDefinition5(t *testing.T) {
	stepsArray := []interface{}{
		map[string]interface{}{
			"name": "step1",
		},
		"invalid",
		map[string]interface{}{
			"name": "step2",
		},
	}
	steps, err := CreateStepsMap(stepsArray)
	if err == nil {
		t.Error(EXPECTED_ERROR_GOT_NIL)
	}
	if steps != nil {
		t.Errorf(EXPECTED_NIL_GOT, steps)
	}
}

func TestCreateStepsMapInvalidStepDefinition6(t *testing.T) {
	stepsArray := []interface{}{
		map[string]interface{}{
			"name": "step1",
		},
		map[string]interface{}{
			"name": "step2",
		},
		"invalid",
	}
	steps, err := CreateStepsMap(stepsArray)
	if err == nil {
		t.Error(EXPECTED_ERROR_GOT_NIL)
	}
	if steps != nil {
		t.Errorf(EXPECTED_NIL_GOT, steps)
	}
}
