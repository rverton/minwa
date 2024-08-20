package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Config struct {
	From  string
	To    string
	Token string
}

func NotifyMail(config Config, subject, msg string) error {
	url := "https://api.postmarkapp.com/email"

	// Create the email payload
	email := map[string]interface{}{
		"From":          config.From,
		"To":            config.To,
		"Subject":       subject,
		"TextBody":      msg,
		"MessageStream": "outbound",
	}

	jsonData, err := json.Marshal(email)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Postmark-Server-Token", config.Token)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	return fmt.Errorf("invalid status returned: %v", resp.StatusCode)
}
