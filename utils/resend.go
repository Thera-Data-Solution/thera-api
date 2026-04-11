package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"thera-api/logger"

	"go.uber.org/zap"
)

type ResendClient struct {
	APIKey string
}

type ResendEmailRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Html    string   `json:"html"`
}

type ResendResponse struct {
	Id string `json:"id"`
}

func NewResendClient(apiKey string) *ResendClient {
	return &ResendClient{APIKey: apiKey}
}

func (c *ResendClient) SendEmail(to, subject, html string) error {
	if c.APIKey == "" {
		return fmt.Errorf("Resend API key tidak ditemukan")
	}

	reqBody := ResendEmailRequest{
		From:    "Theravickya Support <noreply@theravickya.com>",
		To:      []string{to},
		Subject: subject,
		Html:    html,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		logger.Log.Error("Gagal marshal email request", zap.Error(err))
		return fmt.Errorf("gagal membuat request email: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Log.Error("Gagal membuat HTTP request", zap.Error(err))
		return fmt.Errorf("gagal membuat HTTP request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Log.Error("Gagal mengirim email", zap.Error(err))
		return fmt.Errorf("gagal mengirim email: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Error("Gagal membaca response", zap.Error(err))
		return fmt.Errorf("gagal membaca response: %w", err)
	}

	if resp.StatusCode >= 300 {
		logger.Log.Error("Resend API error",
			zap.Int("statusCode", resp.StatusCode),
			zap.String("response", string(body)))
		return fmt.Errorf("Resend API error: %s (status: %d)", string(body), resp.StatusCode)
	}

	var resendResp ResendResponse
	if err := json.Unmarshal(body, &resendResp); err != nil {
		logger.Log.Warn("Gagal unmarshal response", zap.Error(err))
		// Tidak fatal, email mungkin sudah terkirim
	}

	logger.Log.Info("Email berhasil dikirim",
		zap.String("to", to),
		zap.String("emailId", resendResp.Id))

	return nil
}
