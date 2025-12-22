package models

type Setting struct {
	ID             string  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	AppName        string  `json:"appName" gorm:"not null"`
	AppLogo        string  `json:"appLogo" gorm:"not null"`
	AppTitle       string  `json:"appTitle" gorm:"not null"`
	AppDescription *string `json:"appDescription,omitempty"`
	AppTheme       *string `json:"appTheme,omitempty"`

	EnableChatBot        bool `json:"enableChatBot" gorm:"default:false"`
	EnableFacilitator    bool `json:"enableFacilitator" gorm:"default:false"`
	EnablePaymentGateway bool `json:"enablePaymentGateway" gorm:"default:false"`

	TelegramBotToken *string `json:"telegramBotToken,omitempty"`
	TelegramChatId   *string `json:"telegramChatId,omitempty"`

	MailKey    *string `json:"mailKey,omitempty"`
	MailSecret *string `json:"mailSecret,omitempty"`

	DiscordReportId *string `json:"discordReportId,omitempty"`

	Timezone *string `json:"timezone,omitempty"`
	TenantId *string `json:"tenantId,omitempty" gorm:"index"`
}
