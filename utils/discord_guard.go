package utils

import "thera-api/models"

func IsDiscordEnabled(setting *models.Setting) bool {
	if setting == nil {
		return false
	}

	if !setting.EnableChatBot {
		return false
	}

	if setting.DiscordReportId == nil || *setting.DiscordReportId == "" {
		return false
	}

	return true
}


