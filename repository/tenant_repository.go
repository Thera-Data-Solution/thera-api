package repository

import (
	"database/sql"
	"errors"
	"thera-api/models"
)

type TenantRepo struct {
	DB *sql.DB
}

func (r *TenantRepo) FindTenantById(id string) (*models.Tenant, error) {
	query := `SELECT id, name, logo FROM "Tenant" WHERE id=$1 LIMIT 1`
	row := r.DB.QueryRow(query, id)
	var s models.Tenant
	err := row.Scan(&s.ID, &s.Name, &s.Logo)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &s, nil
}
