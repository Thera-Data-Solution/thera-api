package models

import "time"

// Tenant merepresentasikan model Tenant di database.
// Menggunakan struct tag 'gorm:"column:..."' untuk pemetaan ke snake_case di DB jika diperlukan.
type Tenant struct {
	ID        string    `gorm:"primaryKey;column:id"`
	Name      string    `gorm:"column:name"`
	Logo      *string   `gorm:"column:logo"`
	IsActive  bool      `gorm:"column:is_active"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// TableName menentukan nama tabel di database.
func (Tenant) TableName() string {
	return "tenant" // Sesuaikan dengan nama tabel di DB
}
