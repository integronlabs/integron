package helpers

import "testing"

func TestRemoveNull(t *testing.T) {
	input := map[string]interface{}{
		"message": nil,
		"world":   "hello",
	}

	output := RemoveNull(input)

	if output.(map[string]interface{})["message"] != nil && output.(map[string]interface{})["world"] != "hello" {
		t.Errorf("Expected nil, got %v", output)
	}
}

func TestRemoveNullArray(t *testing.T) {
	input := []interface{}{
		map[string]interface{}{
			"message": nil,
		},
		map[string]interface{}{
			"world": "hello",
		},
		nil,
	}

	output := RemoveNull(input)

	if output.([]interface{})[0].(map[string]interface{})["message"] != nil && output.([]interface{})[1].(map[string]interface{})["world"] != "hello" {
		t.Errorf("Expected nil, got %v", output)
	}
}
