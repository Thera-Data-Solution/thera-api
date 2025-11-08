package models

type Categories struct {
	ID             string   `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name           string   `json:"name" gorm:"not null"`
	Description    *string  `json:"description,omitempty"`
	DescriptionEn  *string  `json:"descriptionEn,omitempty"`
	Slug           string   `json:"slug" gorm:"uniqueIndex;not null"`
	Image          *string  `json:"image,omitempty"`
	Start          int      `json:"start"`
	End            int      `json:"end"`
	Location       *string  `json:"location,omitempty"`
	Price          *float64 `json:"price,omitempty"`
	IsGroup        bool     `json:"isGroup" gorm:"default:false"`
	IsFree         bool     `json:"isFree" gorm:"default:false"`
	IsPayAsYouWish bool     `json:"isPayAsYouWish" gorm:"default:false"`
	IsManual       bool     `json:"isManual" gorm:"default:false"`
	Disable        bool     `json:"disable" gorm:"default:false"`
	TenantId       *string  `json:"tenantId,omitempty" gorm:"index"`

	Schedules []Schedules `json:"schedules,omitempty" gorm:"foreignKey:CategoryId"`
}
