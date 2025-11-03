package models

type Setting struct {
	ID                   string  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	AppName              string  `json:"appName" gorm:"not null"`
	AppLogo              string  `json:"appLogo" gorm:"not null"`
	AppTitle             string  `json:"appTitle" gorm:"not null"`
	AppDescription       *string `json:"appDescription,omitempty"`
	AppTheme             *string `json:"appTheme,omitempty"`
	AppMainColor         string  `json:"appMainColor" gorm:"not null"`
	AppHeaderColor       string  `json:"appHeaderColor" gorm:"not null"`
	AppFooterColor       string  `json:"appFooterColor" gorm:"not null"`
	FontSize             int     `json:"fontSize"`
	AppDecoration        *string `json:"appDecoration,omitempty" gorm:"type:text"`
	EnableChatBot        bool    `json:"enableChatBot" gorm:"default:false"`
	EnableFacilitator    bool    `json:"enableFacilitator" gorm:"default:false"`
	EnablePaymentGateway bool    `json:"enablePaymentGateway" gorm:"default:false"`
	MetaOg               *string `json:"metaOg,omitempty" gorm:"type:text"`
	Timezone             *string `json:"timezone,omitempty"`
	TenantId             *string `json:"tenantId,omitempty" gorm:"index"`
}
