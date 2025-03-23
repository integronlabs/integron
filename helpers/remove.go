package helpers

func RemoveNull(input interface{}) interface{} {
	// remove null values from input
	if inputArray, ok := input.([]interface{}); ok {
		for i, val := range inputArray {
			if val == nil {
				inputArray = append(inputArray[:i], inputArray[i+1:]...)
			} else {
				inputArray[i] = RemoveNull(val)
			}
		}
		return inputArray
	}

	if inputMap, ok := input.(map[string]interface{}); ok {
		for key, val := range inputMap {
			if val == nil {
				delete(inputMap, key)
			} else {
				inputMap[key] = RemoveNull(val)
			}
		}
		return inputMap
	}

	return input
}
