package repository

import (
	"database/sql"
	"thera-api/models"
)

type SessionRepository struct {
	DB *sql.DB
}

func (r *SessionRepository) CreateSession(s *models.Session) error {
	query := `INSERT INTO "Session" (id, "userId", "tenantUserId", token, device, ip, "expiresAt", "createdAt", "tenantId")
	VALUES ($1,$2,$3,$4,$5,$6,$7, now(), $8)`
	_, err := r.DB.Exec(query, s.ID, s.UserId, s.TenantUserId, s.Token, s.Device, s.IP, s.ExpiresAt, s.TenantId)
	return err
}

func (r *SessionRepository) DeleteByTenantUserId(tenantUserId string) error {
	query := `DELETE FROM "Session" WHERE "tenantUserId"=$1`
	_, err := r.DB.Exec(query, tenantUserId)
	return err
}

func (r *SessionRepository) FindSessionByToken(token string) (*models.Session, error) {
	query := `SELECT id, "userId", "tenantUserId", token, device, ip, "expiresAt", "tenantId" FROM "Session" WHERE token=$1 LIMIT 1`
	row := r.DB.QueryRow(query, token)
	var s models.Session
	var userId, tenantUserId, device, ip, tenantId sql.NullString
	var expiresAt sql.NullString
	err := row.Scan(&s.ID, &userId, &tenantUserId, &s.Token, &device, &ip, &expiresAt, &tenantId)
	if err != nil {
		return nil, err
	}
	if userId.Valid {
		v := userId.String
		s.UserId = &v
	}
	if tenantUserId.Valid {
		v := tenantUserId.String
		s.TenantUserId = &v
	}
	if device.Valid {
		v := device.String
		s.Device = &v
	}
	if ip.Valid {
		v := ip.String
		s.IP = &v
	}
	if expiresAt.Valid {
		s.ExpiresAt = expiresAt.String
	}
	if tenantId.Valid {
		v := tenantId.String
		s.TenantId = &v
	}
	return &s, nil
}
