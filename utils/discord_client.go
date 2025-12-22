package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type DiscordWebhookConfig struct {
	WebhookURL string
}

type discordWebhookPayload struct {
	Content string `json:"content"`
}

func SendDiscordWebhook(cfg DiscordWebhookConfig, message string) error {
	payload := discordWebhookPayload{
		Content: message,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(cfg.WebhookURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("discord send failed: %s", resp.Status)
	}

	return nil
}


