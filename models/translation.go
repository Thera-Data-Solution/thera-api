package models

type Translation struct {
	ID        string  `json:"id" gorm:"column:id"`
	Locale    string  `json:"locale" gorm:"column:locale"`
	Namespace string  `json:"namespace" gorm:"column:namespace"`
	Key       string  `json:"key" gorm:"column:key"`
	Value     string  `json:"value" gorm:"column:value"`
	TenantId  *string `json:"tenantId,omitempty" gorm:"column:tenantId"`
}

func (Translation) TableName() string {
	return "Translation"
}
