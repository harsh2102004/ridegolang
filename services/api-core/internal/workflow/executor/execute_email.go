package executor

import (
	"encoding/json"
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

func (e *Executor) executeEmail(
	config []byte,
	input []byte,
) ([]byte, error) {

	var cfg struct {
		To      string `json:"to"`
		Subject string `json:"subject"`
		Body    string `json:"body"`
	}

	if err := json.Unmarshal(config, &cfg); err != nil {
		return nil, err
	}

	// Replace placeholders from trigger input
	if len(input) > 0 {
		var inputMap map[string]interface{}

		if err := json.Unmarshal(input, &inputMap); err == nil {

			for key, value := range inputMap {

				placeholder := "{{trigger." + key + "}}"

				valueString := fmt.Sprintf("%v", value)
                fmt.Println("placeholder is "+placeholder)

				cfg.Body = strings.ReplaceAll(
					cfg.Body,
					placeholder,
					valueString,
				)

				cfg.Subject = strings.ReplaceAll(
					cfg.Subject,
					placeholder,
					valueString,
				)
			}
		}
	}

	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	if cfg.To == "" {
		return nil, fmt.Errorf("email recipient is required")
	}

	if cfg.Subject == "" {
		return nil, fmt.Errorf("email subject is required")
	}

	if cfg.Body == "" {
		return nil, fmt.Errorf("email body is required")
	}

	if from == "" || password == "" || smtpHost == "" || smtpPort == "" {
		return nil, fmt.Errorf("smtp configuration is missing")
	}

	message := []byte(
		"Subject: " + cfg.Subject + "\r\n" +
			"\r\n" +
			cfg.Body,
	)

	auth := smtp.PlainAuth(
		"",
		from,
		password,
		smtpHost,
	)

	err := smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		from,
		[]string{cfg.To},
		message,
	)

	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"status":  "sent",
		"to":      cfg.To,
		"subject": cfg.Subject,
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return resultBytes, nil
}
