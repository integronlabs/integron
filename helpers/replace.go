package helpers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/yalp/jsonpath"
)

func Replace(input string, stepOutputs map[string]interface{}) (string, error) {
	re := regexp.MustCompile(`\$\.[a-zA-Z0-9_\[\]\.]+`)
	matches := re.FindAllString(input, -1)

	for _, match := range matches {
		value, err := jsonpath.Read(stepOutputs, match)
		if err != nil {
			return input, err
		}
		input = strings.ReplaceAll(input, match, fmt.Sprintf("%v", value))
	}

	return input, nil
}
