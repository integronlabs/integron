package helpers

import (
	"testing"
)

const HELLO_WORLD = "Hello, world!"
const HELLO_NAME = "Hello, $.name!"

func TestReplace(t *testing.T) {
	// replace json path placeholders with values
	input := HELLO_NAME
	values := map[string]interface{}{
		"name": "world",
	}
	expected := HELLO_WORLD
	result, err := Replace(input, values)
	if err != nil {
		t.Errorf(EXPECTED_NIL_GOT, err)
	}
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestReplaceInvalidJsonPath(t *testing.T) {
	// replace json path placeholders with values
	input := HELLO_NAME
	values := map[string]interface{}{}
	expected := "unknown key name"
	result, err := Replace(input, values)
	if err == nil {
		t.Error(EXPECTED_ERROR_GOT_NIL)
	}
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestReplaceInvalidJsonPath1(t *testing.T) {
	// replace json path placeholders with values
	input := "Hello, $.name[]!"
	values := map[string]interface{}{
		"name": "world",
	}
	expected := "parsing error: $.name[]	:1:8 - 1:9 unexpected \"]\" while scanning extensions"
	result, err := Replace(input, values)
	if err == nil {
		t.Error(EXPECTED_ERROR_GOT_NIL)
	}
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestTransformArray(t *testing.T) {
	// transform array of objects
	input := []interface{}{
		map[string]interface{}{
			"name": "world",
		},
	}
	output := map[string]interface{}{
		"message": HELLO_NAME,
	}
	expected := []map[string]interface{}{
		{
			"message": HELLO_WORLD,
		},
	}
	result, err := TransformArray(input, output)
	if err != nil {
		t.Errorf(EXPECTED_NIL_GOT, err)
	}
	if len(result) != len(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
	for i := range result {
		if result[i].(map[string]interface{})["message"] != expected[i]["message"] {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	}
}

func TestTransformArrayInvalidInputFormat(t *testing.T) {
	// transform array of objects
	input := []interface{}{
		"invalid",
	}
	output := map[string]interface{}{
		"message": HELLO_NAME,
	}
	expected := []map[string]interface{}{}
	result, err := TransformArray(input, output)
	if err == nil {
		t.Error(EXPECTED_ERROR_GOT_NIL)
	}
	if len(result) != len(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
	for i := range result {
		if result[i].(map[string]interface{})["message"] != expected[i]["message"] {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	}
}

func TestTransformArrayInvalidOutputFormat(t *testing.T) {
	// transform array of objects
	input := []interface{}{
		map[string]interface{}{
			"name": "world",
		},
	}
	output := map[string]interface{}{
		"message": "$.phoneNumbers[].type",
	}
	expected := []interface{}{}
	result, err := TransformArray(input, output)
	if err == nil {
		t.Error(EXPECTED_ERROR_GOT_NIL)
	}
	if len(result) != len(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestTransformBody(t *testing.T) {
	// transform object
	input := map[string]interface{}{
		"name": "world",
	}
	output := map[string]interface{}{
		"message": HELLO_NAME,
	}
	expected := map[string]interface{}{
		"message": HELLO_WORLD,
	}
	result, err := TransformBody(input, output)
	if err != nil {
		t.Errorf(EXPECTED_NIL_GOT, err)
	}
	if result.(map[string]interface{})["message"] != expected["message"] {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestTransformBodyInvalidInputFormat(t *testing.T) {
	// transform object
	input := "invalid"
	output := map[string]interface{}{
		"message": HELLO_NAME,
	}
	expected := map[string]interface{}{}
	result, err := TransformBody(input, output)
	if err == nil {
		t.Error(EXPECTED_ERROR_GOT_NIL)
	}
	if result.(map[string]interface{})["message"] != expected["message"] {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestTransformBodyInvalidOutputFormat(t *testing.T) {
	// transform object
	input := map[string]interface{}{
		"name": "world",
	}
	output := map[string]interface{}{
		"message": "$.phoneNumbers[].type",
	}
	expected := map[string]interface{}{}
	result, err := TransformBody(input, output)
	if err == nil {
		t.Error(EXPECTED_ERROR_GOT_NIL)
	}
	if result.(map[string]interface{})["message"] != expected["message"] {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestTransformBodyOutputArray(t *testing.T) {
	// transform object
	input := map[string]interface{}{
		"name": "world",
	}
	output := []interface{}{
		map[string]interface{}{
			"message": HELLO_NAME,
		},
	}
	expected := []interface{}{
		map[string]interface{}{
			"message": HELLO_WORLD,
		},
	}
	result, err := TransformBody(input, output)
	if err != nil {
		t.Errorf(EXPECTED_NIL_GOT, err)
	}
	if result.([]interface{})[0].(map[string]interface{})["message"] != expected[0].(map[string]interface{})["message"] {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestTransformBodyOutputArrayInvalidInputFormat(t *testing.T) {
	// transform object
	input := "invalid"
	output := []interface{}{
		map[string]interface{}{
			"message": HELLO_NAME,
		},
	}
	expected := []interface{}{}
	result, err := TransformBody(input, output)
	if err == nil {
		t.Error(EXPECTED_ERROR_GOT_NIL)
	}
	if len(result.([]interface{})) != len(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestTransformBodyOutputArrayInvalidOutputFormat(t *testing.T) {
	// transform object
	input := map[string]interface{}{
		"name": "world",
	}
	output := []interface{}{
		map[string]interface{}{
			"message": "$.phoneNumbers[].type",
		},
	}
	expected := []interface{}{}
	result, err := TransformBody(input, output)
	if err == nil {
		t.Error(EXPECTED_ERROR_GOT_NIL)
	}
	if len(result.([]interface{})) != len(expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestTransformBodyOutputInteger(t *testing.T) {
	// transform object
	input := map[string]interface{}{
		"name": "world",
	}
	output := 42
	expected := 42
	result, err := TransformBody(input, output)
	if err != nil {
		t.Errorf(EXPECTED_NIL_GOT, err)
	}
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
