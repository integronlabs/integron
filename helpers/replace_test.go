package helpers

import (
	"testing"
)

const HELLO_WORLD = "Hello, world!"
const HELLO_NAME = "Hello, $.name!"

const EXPECTED_BUT_GOT = "Expected %v, got %v"
const INVALID_JSON_PATH = "$.phoneNumbers[].type"

func TestReplace(t *testing.T) {
	// replace json path placeholders with values
	input := HELLO_NAME
	values := map[string]interface{}{
		"name": "world",
	}
	expected := HELLO_WORLD
	result := Replace(input, values)
	if result != expected {
		t.Errorf(EXPECTED_BUT_GOT, expected, result)
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
	result := TransformArray(input, output)
	if len(result) != len(expected) {
		t.Errorf(EXPECTED_BUT_GOT, expected, result)
	}
	for i := range result {
		if result[i].(map[string]interface{})["message"] != expected[i]["message"] {
			t.Errorf(EXPECTED_BUT_GOT, expected, result)
		}
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
	result := TransformBody(input, output)
	if result.(map[string]interface{})["message"] != expected["message"] {
		t.Errorf(EXPECTED_BUT_GOT, expected, result)
	}
}

func TestTransformBodyJsonPath(t *testing.T) {
	// transform object
	input := map[string]interface{}{
		"name": "world",
	}
	output := map[string]interface{}{
		"message": "$.name",
	}
	expected := map[string]interface{}{
		"message": "world",
	}
	result := TransformBody(input, output)
	if result.(map[string]interface{})["message"] != expected["message"] {
		t.Errorf(EXPECTED_BUT_GOT, expected, result)
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
	result := TransformBody(input, output)
	if result.([]interface{})[0].(map[string]interface{})["message"] != expected[0].(map[string]interface{})["message"] {
		t.Errorf(EXPECTED_BUT_GOT, expected, result)
	}
}

func TestTransformBodyOutputInteger(t *testing.T) {
	// transform object
	input := map[string]interface{}{
		"name": "world",
	}
	output := 42
	expected := 42
	result := TransformBody(input, output)
	if result != expected {
		t.Errorf(EXPECTED_BUT_GOT, expected, result)
	}
}
