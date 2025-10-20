package controllers

import (
	"thera-api/repository"
)

type AuthController struct {
	UserRepo       *repository.UserRepository
	TenantUserRepo *repository.TenantUserRepository
	SessionRepo    *repository.SessionRepository
}

// ---- User register/login ----
// func (ac *AuthController) UserRegister(c echo.Context) error {
// 	type req struct {
// 		Email    string  `json:"email"`
// 		Password string  `json:"password"`
// 		FullName string  `json:"fullName"`
// 		Phone    *string `json:"phone"`
// 		TenantId *string `json:"tenantId"` // optional
// 	}
// 	var r req
// 	if err := c.Bind(&r); err != nil {
// 		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
// 	}

// 	// check existing
// 	_, err := ac.UserRepo.FindByEmailAndTenant(r.Email, r.TenantId)
// 	if err == nil {
// 		return c.JSON(http.StatusBadRequest, echo.Map{"error": "email already used"})
// 	}

// 	hashed, _ := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
// 	avatar := "https://api.dicebear.com/9.x/fun-emoji/svg?seed=" + url.QueryEscape(r.FullName)
// 	u := &models.User{
// 		ID:       uuid.New().String(),
// 		Email:    r.Email,
// 		Password: string(hashed),
// 		FullName: r.FullName,
// 		Phone:    r.Phone,
// 		TenantId: r.TenantId,
// 		Avatar:   &avatar,
// 	}

// 	if err := ac.UserRepo.CreateUser(u); err != nil {
// 		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
// 	}

// 	return c.JSON(http.StatusOK, echo.Map{"message": "register success"})
// }

// func (ac *AuthController) UserLogin(c echo.Context) error {
// 	type req struct {
// 		Email    string  `json:"email"`
// 		Password string  `json:"password"`
// 		TenantId *string `json:"tenantId"`
// 		Device   *string `json:"device"`
// 		IP       *string `json:"ip"`
// 	}
// 	var r req
// 	if err := c.Bind(&r); err != nil {
// 		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
// 	}

// 	user, err := ac.UserRepo.FindByEmailAndTenant(r.Email, r.TenantId)
// 	if err != nil {
// 		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "email/password wrong"})
// 	}

// 	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(r.Password)) != nil {
// 		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "email/password wrong"})
// 	}

// 	// create session
// 	token := uuid.New().String()
// 	expHours := 24
// 	if v := os.Getenv("SESSION_EXPIRE_HOURS"); v != "" {
// 		if h, err := strconv.Atoi(v); err == nil {
// 			expHours = h
// 		}
// 	}
// 	exp := time.Now().Add(time.Hour * time.Duration(expHours)).UTC()

// 	s := &models.Session{
// 		ID:        uuid.New().String(),
// 		UserId:    &user.ID,
// 		Token:     token,
// 		Device:    r.Device,
// 		IP:        r.IP,
// 		ExpiresAt: exp.Format(time.RFC3339),
// 		TenantId:  r.TenantId,
// 	}
// 	if err := ac.SessionRepo.CreateSession(s); err != nil {
// 		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "gagal membuat session"})
// 	}

// 	// return token
// 	return c.JSON(http.StatusOK, echo.Map{
// 		"token":     token,
// 		"expiresAt": s.ExpiresAt,
// 	})
// }
