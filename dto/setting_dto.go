package dto

type SettingRequestBody struct {
	AppName        string  `json:"appName" form:"appName" gorm:"not null"`
	AppLogo        string  `json:"appLogo" form:"appLogo"`
	AppTitle       string  `json:"appTitle" form:"appTitle" gorm:"not null"`
	AppDescription *string `json:"appDescription,omitempty" form:"appDescription"`
	AppTheme       *string `json:"appTheme,omitempty" form:"appTheme"`

	EnableChatBot        bool `json:"enableChatBot" form:"enableChatBot" gorm:"default:false"`
	EnableFacilitator    bool `json:"enableFacilitator" form:"enableFacilitator" gorm:"default:false"`
	EnablePaymentGateway bool `json:"enablePaymentGateway" form:"enablePaymentGateway" gorm:"default:false"`

	Timezone *string `json:"timezone,omitempty" form:"timezone"`

	// 🔹 Telegram
	TelegramBotToken *string `json:"telegramBotToken,omitempty" form:"telegramBotToken"`
	TelegramChatId   *string `json:"telegramChatId,omitempty" form:"telegramChatId"`

	// 🔹 Email
	MailKey    *string `json:"mailKey,omitempty" form:"mailKey"`
	MailSecret *string `json:"mailSecret,omitempty" form:"mailSecret"`

	DiscordReportId *string `json:"discordReportId,omitempty" form:"discordReportId"`

	TenantId *string `json:"tenantId,omitempty" form:"tenantId" gorm:"index"`
}

type SettingGetResponse struct {
	AppName        string  `json:"appName"`
	AppLogo        string  `json:"appLogo"`
	AppTitle       string  `json:"appTitle"`
	AppDescription *string `json:"appDescription,omitempty"`
	AppTheme       *string `json:"appTheme,omitempty"`

	Timezone *string `json:"timezone,omitempty"`
}

type SettingResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
