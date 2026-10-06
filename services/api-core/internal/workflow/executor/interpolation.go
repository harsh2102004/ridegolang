package executor

import (
	"encoding/json"
	"fmt"
	"strings"
)

func interpolateString(template string, input []byte) string {
	if len(input) == 0 {
		return template
	}

	var data map[string]interface{}
	if err := json.Unmarshal(input, &data); err != nil {
		return template
	}

	result := template

	for key, value := range data {
		valueString := fmt.Sprintf("%v", value)

		if str, ok := value.(string); ok {
			valueString = str
		}

		result = strings.ReplaceAll(result, "{{"+key+"}}", valueString)
		result = strings.ReplaceAll(result, "{{trigger."+key+"}}", valueString)
		result = strings.ReplaceAll(result, "{{input."+key+"}}", valueString)
	}

	return result
}