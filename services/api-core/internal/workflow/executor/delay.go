package executor

import (
	"encoding/json"
	"time"
)

func (e *Executor) executeDelay(
	config []byte,
	input []byte,
) ([]byte, error) {

	var cfg struct {
		Duration int    `json:"duration"`
		Unit     string `json:"unit"`
	}

	if err := json.Unmarshal(config, &cfg); err != nil {
		return nil, err
	}

	if cfg.Duration <= 0 {
		cfg.Duration = 1
	}

	var wait time.Duration

	switch cfg.Unit {
	case "minute":
		wait = time.Duration(cfg.Duration) * time.Minute

	case "hour":
		wait = time.Duration(cfg.Duration) * time.Hour

	case "day":
		wait = time.Duration(cfg.Duration) * 24 * time.Hour

	default:
		wait = time.Duration(cfg.Duration) * time.Second
	}

	time.Sleep(wait)

	// preserve original payload
	return input, nil
}