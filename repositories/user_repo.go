package repositories

import (
	"thera-api/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func (r *UserRepository) FindByEmailAndTenant(email, tenantId string) (*models.User, error) {
	var user models.User
	err := r.DB.Where(`email = ? AND tenant_id = ?`, email, tenantId).First(&user).Error
	return &user, err
}

func (r *UserRepository) FindByEmailAndTenantForLogin(email, tenantId string) (*models.User, error) {
	var user models.User
	err := r.DB.Where(`email = ? AND tenant_id = ? AND disable = false`, email, tenantId).First(&user).Error
	return &user, err
}

func (r *UserRepository) FindByPhoneAndTenant(phone, tenantId string) (*models.User, error) {
	var user models.User
	err := r.DB.Where(`phone = ? AND tenant_id = ?`, phone, tenantId).First(&user).Error
	return &user, err
}

func (r *UserRepository) Create(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) FindByID(idToken string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("id = ?", idToken).First(&user).Error
	return &user, err
}

func (r *UserRepository) FindAllByTenantId(tenantId string) ([]models.User, error) {
	var users []models.User
	err := r.DB.Where("tenant_id = ?", tenantId).Find(&users).Error
	return users, err
}

func (r *UserRepository) GetAllByTenantPaginated(tenantId string, page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// Hitung total
	if err := r.DB.
		Model(&models.User{}).
		Where("tenant_id = ?", tenantId).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	err := r.DB.
		Where("tenant_id = ?", tenantId).
		Limit(pageSize).
		Offset(offset).
		Find(&users).Error

	if err != nil {
		return nil, 0, err
	}

	// Hapus password (jika masih ada field password)
	for i := range users {
		users[i].Password = ""
	}

	return users, total, nil
}

func (r *UserRepository) GetAllPaginated(page, pageSize int) ([]models.User, []string, int64, error) {
	var users []models.User
	var tenantNames []string
	var total int64

	// Hitung total
	if err := r.DB.
		Model(&models.User{}).
		Count(&total).Error; err != nil {
		return nil, nil, 0, err
	}

	offset := (page - 1) * pageSize

	// Query dengan JOIN
	rows, err := r.DB.
		Table("users u").
		Select(`
            u.id,
            u.avatar,
            u.email,
            u.full_name,
            u.phone,
            u.address,
            u.ig,
            u.fb,
            u.disable,
            u.tenant_id,
            t.name AS tenant_name
        `).
		Joins("LEFT JOIN tenants t ON t.id::text = u.tenant_id").
		Limit(pageSize).
		Offset(offset).
		Rows()
	if err != nil {
		return nil, nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var u models.User
		var tenantName string

		err := rows.Scan(
			&u.ID,
			&u.Avatar,
			&u.Email,
			&u.FullName,
			&u.Phone,
			&u.Address,
			&u.Ig,
			&u.Fb,
			&u.Disable,
			&u.TenantId,
			&tenantName,
		)
		if err != nil {
			return nil, nil, 0, err
		}

		u.Password = "" // biar aman
		users = append(users, u)
		tenantNames = append(tenantNames, tenantName)
	}

	return users, tenantNames, total, nil
}

func (r *UserRepository) Update(user *models.User) error {
	return r.DB.Save(user).Error
}
