package models

type Translation struct {
	ID        string  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Locale    string  `json:"locale" gorm:"not null;index:idx_locale_namespace_key,unique"`
	Namespace string  `json:"namespace" gorm:"not null;index:idx_locale_namespace_key,unique"`
	Key       string  `json:"key" gorm:"not null;index:idx_locale_namespace_key,unique"`
	Value     string  `json:"value" gorm:"not null"`
	TenantId  *string `json:"tenantId,omitempty" gorm:"index"`
}

// TranslationFilter digunakan untuk advanced filtering di endpoint GET /translations
type TranslationFilter struct {
	Locale    *string `json:"locale"`
	Namespace *string `json:"namespace"`
	Key       *string `json:"key"`
	Value     *string `json:"value"`
	Search    *string `json:"search"` // pencarian bebas di key / value / namespace
}

// PaginationRequest menyimpan informasi pagination dari client
type PaginationRequest struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// PaginatedTranslations adalah bentuk response dengan metadata pagination
type PaginatedTranslations struct {
	Data       []Translation `json:"data"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
	TotalItems int64         `json:"totalItems"`
	TotalPages int           `json:"totalPages"`
}

