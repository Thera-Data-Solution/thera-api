package utils

import (
	"fmt"
	"net/http"
	"net/url"
)

type TelegramConfig struct {
	BotToken string
	ChatID   string
}

func SendTelegramHTML(cfg TelegramConfig, message string) error {
	endpoint := fmt.Sprintf(
		"https://api.telegram.org/bot%s/sendMessage",
		cfg.BotToken,
	)

	data := url.Values{}
	data.Set("chat_id", cfg.ChatID)
	data.Set("text", message)
	data.Set("parse_mode", "HTML")
	data.Set("disable_web_page_preview", "true")

	resp, err := http.PostForm(endpoint, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("telegram send failed: %s", resp.Status)
	}

	return nil
}
