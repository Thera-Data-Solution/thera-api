package utils

import (
	"thera-api/logger"
	"thera-api/models"
)

func IsTelegramEnabled(setting *models.Setting) bool {
	if setting == nil {
		logger.Log.Warn("Setting is nil")
		return false
	}

	if !setting.EnableChatBot {
		logger.Log.Warn("Chat bot is disabled")
		return false
	}

	if setting.TelegramBotToken == nil || *setting.TelegramBotToken == "" {
		logger.Log.Warn("Telegram bot token is missing")
		return false
	}

	if setting.TelegramChatId == nil || *setting.TelegramChatId == "" {
		logger.Log.Warn("Telegram chat ID is missing")
		return false
	}

	return true
}
