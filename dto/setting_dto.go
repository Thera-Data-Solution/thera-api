package dto

type SettingRequestBody struct {
	AppName              string  `json:"appName" form:"appName" gorm:"not null"`
	AppLogo              string  `json:"appLogo" form:"appLogo" gorm:"not null"`
	AppTitle             string  `json:"appTitle" form:"appTitle" gorm:"not null"`
	AppDescription       *string `json:"appDescription,omitempty" form:"appDescription"`
	AppTheme             *string `json:"appTheme,omitempty" form:"appTheme"`
	AppMainColor         string  `json:"appMainColor" form:"appMainColor" gorm:"not null"`
	AppHeaderColor       string  `json:"appHeaderColor" form:"appHeaderColor" gorm:"not null"`
	AppFooterColor       string  `json:"appFooterColor" form:"appFooterColor" gorm:"not null"`
	FontSize             int     `json:"fontSize" form:"fontSize"`
	AppDecoration        *string `json:"appDecoration,omitempty" form:"appDecoration" gorm:"type:text"`
	EnableChatBot        bool    `json:"enableChatBot" form:"enableChatBot" gorm:"default:false"`
	EnableFacilitator    bool    `json:"enableFacilitator" form:"enableFacilitator" gorm:"default:false"`
	EnablePaymentGateway bool    `json:"enablePaymentGateway" form:"enablePaymentGateway" gorm:"default:false"`
	MetaOg               *string `json:"metaOg,omitempty" form:"metaOg" gorm:"type:text"`
	Timezone             *string `json:"timezone,omitempty" form:"timezone"`
	TenantId             *string `json:"tenantId,omitempty" form:"tenantId" gorm:"index"`
}