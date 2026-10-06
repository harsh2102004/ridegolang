package executor

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

func (e *Executor) executeHTTPRequest(config []byte, input []byte) ([]byte, error) {
	var cfg struct {
		Method         string          `json:"method"`
		URL            string          `json:"url"`
		Body           json.RawMessage `json:"body"`
		TimeoutSeconds int             `json:"timeout_seconds"`
	}

	if err := json.Unmarshal(config, &cfg); err != nil {
		return nil, err
	}

	cfg.URL = interpolateString(cfg.URL, input)

	if len(input) > 0 && len(cfg.Body) > 0 {
		bodyString := string(cfg.Body)
		bodyString = interpolateString(bodyString, input)
		cfg.Body = json.RawMessage(bodyString)
	}

	if cfg.URL == "" {
		return nil, errors.New("http request url is required")
	}

	if cfg.Method == "" {
		cfg.Method = "GET"
	}

	var bodyBytes []byte

	if len(cfg.Body) > 0 && string(cfg.Body) != "null" {
		var bodyString string

		if err := json.Unmarshal(cfg.Body, &bodyString); err == nil {
			bodyBytes = []byte(bodyString)
		} else {
			bodyBytes = cfg.Body
		}
	}

	req, err := http.NewRequest(
		strings.ToUpper(cfg.Method),
		cfg.URL,
		bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	if cfg.TimeoutSeconds <= 0 {
		cfg.TimeoutSeconds = 15
	}

	client := &http.Client{
		Timeout: time.Duration(cfg.TimeoutSeconds) * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var output map[string]interface{}

	if len(input) > 0 {
		if err := json.Unmarshal(input, &output); err != nil {
			output = map[string]interface{}{}
		}
	} else {
		output = map[string]interface{}{}
	}

	output["http"] = map[string]interface{}{
		"status_code": resp.StatusCode,
		"status":      resp.Status,
		"body":        string(responseBody),
	}

	resultBytes, err := json.Marshal(output)
	if err != nil {
		return nil, err
	}

	return resultBytes, nil
}
