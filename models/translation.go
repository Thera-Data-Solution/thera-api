package models

type Translation struct {
	ID        string  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Locale    string  `json:"locale" gorm:"not null;index:idx_locale_namespace_key,unique"`
	Namespace string  `json:"namespace" gorm:"not null;index:idx_locale_namespace_key,unique"`
	Key       string  `json:"key" gorm:"not null;index:idx_locale_namespace_key,unique"`
	Value     string  `json:"value" gorm:"not null"`
	TenantId  *string `json:"tenantId,omitempty" gorm:"index"`
}
