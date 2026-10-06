package executor

import (
	"encoding/json"
	"fmt"
)

func (e *Executor) executeCondition(config []byte, input []byte) ([]byte, error) {
	var cfg struct {
		Field    string      `json:"field"`
		Operator string      `json:"operator"`
		Value    interface{} `json:"value"`
	}

	if err := json.Unmarshal(config, &cfg); err != nil {
		return nil, err
	}

	if cfg.Field == "" {
		return nil, fmt.Errorf("condition field is required")
	}

	if cfg.Operator == "" {
		return nil, fmt.Errorf("condition operator is required")
	}

	var inputMap map[string]interface{}
	if err := json.Unmarshal(input, &inputMap); err != nil {
		return nil, err
	}

	actualValue := inputMap[cfg.Field]

	result := false

	switch cfg.Operator {
	case "equals":
		result = actualValue == cfg.Value
	case "not_equals":
		result = actualValue != cfg.Value
	default:
		return nil, fmt.Errorf("unsupported condition operator: %s", cfg.Operator)
	}

	inputMap["condition_result"] = result

	return json.Marshal(inputMap)
}
