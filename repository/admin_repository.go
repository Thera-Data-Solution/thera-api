package repository

import (
	"database/sql"
	"errors"
	"thera-api/models"
)

type TenantUserRepository struct {
	DB *sql.DB
}

func (r *TenantUserRepository) CreateTenantUser(t *models.TenantUser) error {
	query := `INSERT INTO "TenantUser" (id, avatar, email, password, "fullName", role, "tenantId")
	VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := r.DB.Exec(query, t.ID, t.Avatar, t.Email, t.Password, t.FullName, "USER", t.TenantId)
	return err
}

func (r *TenantUserRepository) FindTenantUserByEmailAndTenant(email string, tenantId *string) (*models.TenantUser, error) {
	query := `SELECT * FROM "TenantUser" WHERE email=$1 AND "tenantId" = $2 LIMIT 1`
	row := r.DB.QueryRow(query, email, tenantId)
	var t models.TenantUser
	err := row.Scan(&t.ID, &t.Avatar, &t.Email, &t.Password, &t.FullName, &t.Role, &t.TenantId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &t, nil
}

func (r *TenantUserRepository) FindTenantUserByEmail(email string) (*models.TenantUser, error) {
	query := `
		SELECT id, "fullName", email, "password", role, "tenantId", avatar
		FROM "TenantUser"
		WHERE LOWER(email) = LOWER($1)
		LIMIT 1
	`
	row := r.DB.QueryRow(query, email)

	var t models.TenantUser
	err := row.Scan(&t.ID, &t.FullName, &t.Email, &t.Password, &t.Role, &t.TenantId, &t.Avatar)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // biar aman, bukan error fatal
		}
		return nil, err
	}

	return &t, nil
}

func (r *TenantUserRepository) FindTenantUserById(id string) (*models.TenantUser, error) {
	query := `SELECT "id", "fullName", "role", "email", "avatar" FROM "TenantUser" WHERE id=$1 LIMIT 1`
	row := r.DB.QueryRow(query, id)
	var u models.TenantUser
	err := row.Scan(&u.ID, &u.FullName, &u.Role, &u.Email, &u.Avatar)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &u, nil
}
