package executor

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"org/api-core/internal/db"

	"github.com/google/uuid"
)


func NewExecutor(geminiService GeminiService) *Executor {
	return &Executor{
        geminiService: geminiService,
    }
}

type RetryConfig struct {
	Enabled      bool `json:"enabled"`
	MaxAttempts  int  `json:"max_attempts"`
	DelaySeconds int  `json:"delay_seconds"`
}

type StepConfig struct {
	Retry RetryConfig `json:"retry"`
}

func (e *Executor) ExecuteStep(
	userID uuid.UUID,
	step db.WorkflowStep,
	input []byte,
) ([]byte, error) {
	retryConfig := getRetryConfig(step.Config)

	if !retryConfig.Enabled {
		return e.executeStepOnce(
			userID,
			step,
			input,
		)
	}

	if retryConfig.MaxAttempts <= 0 {
		retryConfig.MaxAttempts = 1
	}

	if retryConfig.DelaySeconds < 0 {
		retryConfig.DelaySeconds = 0
	}

	var output []byte
	var err error

	for attempt := 1; attempt <= retryConfig.MaxAttempts; attempt++ {
		output, err = e.executeStepOnce(
			userID,
			step,
			input,
		)

		if err == nil {
			return output, nil
		}

		if attempt < retryConfig.MaxAttempts {
			time.Sleep(time.Duration(retryConfig.DelaySeconds) * time.Second)
		}
	}

	return nil, err
}

func (e *Executor) executeStepOnce(
	userID uuid.UUID,
	step db.WorkflowStep,
	input []byte,
) ([]byte, error) {
	switch step.StepType {

	case "webhook_trigger":

		if len(input) > 0 {
			return input, nil
		}

		return step.Config, nil

	case "http_request":
		return e.executeHTTPRequest(step.Config, input)

	case "condition":
		return e.executeCondition(step.Config, input)

	case "delay":
		return e.executeDelay(step.Config, input)

	case "email":
		return e.executeEmail(step.Config, input)

	case "ai":
		return e.executeAI(
            context.Background(),
			userID,
			step.Config,
			input,
		)

	default:
		return nil, errors.New(
			"unsupported step type: " + step.StepType,
		)
	}
}

func getRetryConfig(config []byte) RetryConfig {
	var stepConfig StepConfig

	err := json.Unmarshal(config, &stepConfig)
	if err != nil {
		return RetryConfig{
			Enabled:      false,
			MaxAttempts:  1,
			DelaySeconds: 0,
		}
	}

	if stepConfig.Retry.MaxAttempts <= 0 {
		stepConfig.Retry.MaxAttempts = 1
	}

	if stepConfig.Retry.DelaySeconds < 0 {
		stepConfig.Retry.DelaySeconds = 0
	}

	return stepConfig.Retry
}
