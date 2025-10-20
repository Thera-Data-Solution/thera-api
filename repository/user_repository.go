package repository

import (
	"database/sql"
	"errors"
	"thera-api/models"
)

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) CreateUser(user *models.User) error {
	query := `INSERT INTO "User" (id, avatar, email, password, "fullName", phone, address, ig, fb, disable, "tenantId")
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	_, err := r.DB.Exec(query,
		user.ID, user.Avatar, user.Email, user.Password, user.FullName,
		user.Phone, user.Address, user.IG, user.FB, user.Disable, user.TenantId,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) FindByEmailAndTenant(email string, tenantId *string) (*models.User, error) {
	query := `SELECT id, avatar, email, password, "fullName", phone, address, ig, fb, disable, "tenantId" FROM "User" WHERE email=$1 AND "tenantId" IS NOT DISTINCT FROM $2 LIMIT 1`
	row := r.DB.QueryRow(query, email, tenantId)
	var u models.User
	var avatar, phone, address, ig, fb, tenant sql.NullString
	err := row.Scan(&u.ID, &avatar, &u.Email, &u.Password, &u.FullName, &phone, &address, &ig, &fb, &u.Disable, &tenant)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	if avatar.Valid {
		v := avatar.String
		u.Avatar = &v
	}
	if phone.Valid {
		v := phone.String
		u.Phone = &v
	}
	if address.Valid {
		v := address.String
		u.Address = &v
	}
	if ig.Valid {
		v := ig.String
		u.IG = &v
	}
	if fb.Valid {
		v := fb.String
		u.FB = &v
	}
	if tenant.Valid {
		v := tenant.String
		u.TenantId = &v
	}
	return &u, nil
}

func (r *UserRepository) FindById(id string) (*models.User, error) {
	query := `SELECT * FROM "User" WHERE id=$1 LIMIT 1`
	row := r.DB.QueryRow(query, id)
	var u models.User
	var avatar, phone, address, ig, fb, tenant sql.NullString
	err := row.Scan(&u.ID, &avatar, &u.Email, &u.FullName, &phone, &address, &ig, &fb, &u.Disable, &tenant)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &u, nil
}
