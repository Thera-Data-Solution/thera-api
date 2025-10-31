package models

type Setting struct {
	ID                   string  `json:"id" gorm:"column:id"`
	AppName              string  `json:"appName" gorm:"column:appName"`
	AppLogo              string  `json:"appLogo" gorm:"column:appLogo"`
	AppTitle             string  `json:"appTitle" gorm:"column:appTitle"`
	AppDescription       *string `json:"appDescription,omitempty" gorm:"column:appDescription"`
	AppTheme             *string `json:"appTheme,omitempty" gorm:"column:appTheme"`
	AppMainColor         string  `json:"appMainColor" gorm:"column:appMainColor"`
	AppHeaderColor       string  `json:"appHeaderColor" gorm:"column:appHeaderColor"`
	AppFooterColor       string  `json:"appFooterColor" gorm:"column:appFooterColor"`
	FontSize             int     `json:"fontSize" gorm:"column:fontSize"`
	AppDecoration        *string `json:"appDecoration,omitempty" gorm:"column:appDecoration;type:json"`
	EnableChatBot        bool    `json:"enableChatBot" gorm:"column:enableChatBot"`
	EnableFacilitator    bool    `json:"enableFacilitator" gorm:"column:enableFacilitator"`
	EnablePaymentGateway bool    `json:"enablePaymentGateway" gorm:"column:enablePaymentGateway"`
	MetaOg               *string `json:"metaOg,omitempty" gorm:"column:metaOg;type:json"`
	Timezone             *string `json:"timezone,omitempty" gorm:"column:timezone"`
	TenantId             *string `json:"tenantId,omitempty" gorm:"column:tenantId"`
}

func (Setting) TableName() string {
	return "Setting"
}
