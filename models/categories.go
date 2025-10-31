package models

type Categories struct {
	ID             string   `json:"id" gorm:"column:id"`
	Name           string   `json:"name" gorm:"column:name"`
	Description    *string  `json:"description,omitempty" gorm:"column:description"`
	DescriptionEn  *string  `json:"descriptionEn,omitempty" gorm:"column:descriptionEn"`
	Slug           string   `json:"slug" gorm:"column:slug"`
	Image          *string  `json:"image,omitempty" gorm:"column:image"`
	Start          int      `json:"start" gorm:"column:start"`
	End            int      `json:"end" gorm:"column:end"`
	Location       *string  `json:"location,omitempty" gorm:"column:location"`
	Price          *float64 `json:"price,omitempty" gorm:"column:price"`
	IsGroup        bool     `json:"isGroup" gorm:"column:isGroup"`
	IsFree         bool     `json:"isFree" gorm:"column:isFree"`
	IsPayAsYouWish bool     `json:"isPayAsYouWish" gorm:"column:isPayAsYouWish"`
	IsManual       bool     `json:"isManual" gorm:"column:isManual"`
	Disable        bool     `json:"disable" gorm:"column:disable"`
	TenantId       *string  `json:"tenantId,omitempty" gorm:"column:tenantId"`
}

func (Categories) TableName() string {
	return "Categories"
}
