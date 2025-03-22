package helpers

import (
	"testing"
)

func TestReplace(t *testing.T) {
	// replace json path placeholders with values
	input := "Hello, $.name!"
	values := map[string]interface{}{
		"name": "world",
	}
	expected := "Hello, world!"
	result := Replace(input, values)
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestReplaceInvalidJsonPath(t *testing.T) {
	// replace json path placeholders with values
	input := "Hello, $.name!"
	values := map[string]interface{}{}
	expected := "Hello, <nil>!"
	result := Replace(input, values)
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
		"message": "Hello, $.name!",
	}
	expected := []map[string]interface{}{
		{
			"message": "Hello, world!",
		},
	}
	result, err := TransformArray(input, output)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
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
		"message": "Hello, $.name!",
	}
	expected := []map[string]interface{}{
		{
			"message": "Hello, <nil>!",
		},
	}
	result, err := TransformArray(input, output)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
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
		t.Errorf("Expected error, got nil")
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
		"message": "Hello, $.name!",
	}
	expected := map[string]interface{}{
		"message": "Hello, world!",
	}
	result, err := TransformBody(input, output)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result.(map[string]interface{})["message"] != expected["message"] {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestTransformBodyInvalidInputFormat(t *testing.T) {
	// transform object
	input := "invalid"
	output := map[string]interface{}{
		"message": "Hello, $.name!",
	}
	expected := map[string]interface{}{
		"message": "Hello, <nil>!",
	}
	result, err := TransformBody(input, output)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
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
		t.Errorf("Expected error, got nil")
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
			"message": "Hello, $.name!",
		},
	}
	expected := []interface{}{
		map[string]interface{}{
			"message": "Hello, world!",
		},
	}
	result, err := TransformBody(input, output)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
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
			"message": "Hello, $.name!",
		},
	}
	expected := []interface{}{
		map[string]interface{}{
			"message": "Hello, <nil>!",
		},
	}
	result, err := TransformBody(input, output)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result.([]interface{})[0].(map[string]interface{})["message"] != expected[0].(map[string]interface{})["message"] {
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
		t.Errorf("Expected error, got nil")
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
		t.Errorf("Expected no error, got %v", err)
	}
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
