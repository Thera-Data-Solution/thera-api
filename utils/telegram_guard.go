package utils

import "thera-api/models"

func IsTelegramEnabled(setting *models.Setting) bool {
	if setting == nil {
		return false
	}

	if !setting.EnableChatBot {
		return false
	}

	if setting.TelegramBotToken == nil || *setting.TelegramBotToken == "" {
		return false
	}

	if setting.TelegramChatId == nil || *setting.TelegramChatId == "" {
		return false
	}

	return true
}
