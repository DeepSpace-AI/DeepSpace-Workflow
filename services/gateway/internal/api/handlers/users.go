package handlers

import (
	"net/http"
	"strconv"

	"deepspace/internal/service/auth"
	"deepspace/internal/service/user"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userSvc *user.Service
	authSvc *auth.UserAuthService
}

func NewUserHandler(userSvc *user.Service, authSvc *auth.UserAuthService) *UserHandler {
	return &UserHandler{userSvc: userSvc, authSvc: authSvc}
}

type updateUserProfileRequest struct {
	DisplayName *string `json:"display_name"`
	FullName    *string `json:"full_name"`
	Title       *string `json:"title"`
	AvatarURL   *string `json:"avatar_url"`
	Bio         *string `json:"bio"`
	Phone       *string `json:"phone"`
}

type updateUserSettingsRequest struct {
	Theme    *string `json:"theme"`
	Locale   *string `json:"locale"`
	Timezone *string `json:"timezone"`
}

type updateUserRequest struct {
	Profile  *updateUserProfileRequest  `json:"profile"`
	Settings *updateUserSettingsRequest `json:"settings"`
}

type changePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		respondInternal(c, "user_id missing")
		return
	}

	userModel, profile, settings, err := h.userSvc.GetMe(c.Request.Context(), userID)
	if err != nil {
		if err == user.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		respondInternal(c, "failed to load user")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":            userModel.ID,
			"email":         userModel.Email,
			"status":        userModel.Status,
			"role":          userModel.Role,
			"last_login_at": userModel.LastLoginAt,
			"created_at":    userModel.CreatedAt,
		},
		"profile":  profile,
		"settings": settings,
	})
}

func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		respondInternal(c, "user_id missing")
		return
	}

	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	var profileUpdate *user.UpdateProfile
	if req.Profile != nil {
		profileUpdate = &user.UpdateProfile{
			DisplayName: req.Profile.DisplayName,
			FullName:    req.Profile.FullName,
			Title:       req.Profile.Title,
			AvatarURL:   req.Profile.AvatarURL,
			Bio:         req.Profile.Bio,
			Phone:       req.Profile.Phone,
		}
	}

	var settingsUpdate *user.UpdateSettings
	if req.Settings != nil {
		settingsUpdate = &user.UpdateSettings{
			Theme:    req.Settings.Theme,
			Locale:   req.Settings.Locale,
			Timezone: req.Settings.Timezone,
		}
	}

	userModel, profile, settings, err := h.userSvc.UpdateMe(c.Request.Context(), userID, profileUpdate, settingsUpdate)
	if err != nil {
		switch err {
		case user.ErrInvalidSetting:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid settings"})
			return
		case user.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		default:
			respondInternal(c, "failed to update user")
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":            userModel.ID,
			"email":         userModel.Email,
			"status":        userModel.Status,
			"role":          userModel.Role,
			"last_login_at": userModel.LastLoginAt,
			"created_at":    userModel.CreatedAt,
		},
		"profile":  profile,
		"settings": settings,
	})
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		respondInternal(c, "user_id missing")
		return
	}

	var req changePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if len(req.NewPassword) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password too short"})
		return
	}

	if err := h.authSvc.ChangePassword(c.Request.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		if err == auth.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		respondInternal(c, "failed to change password")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func getUserID(c *gin.Context) (int64, bool) {
	value, ok := c.Get("user_id")
	if !ok {
		return 0, false
	}

	id, ok := value.(int64)
	if !ok {
		return 0, false
	}
	return id, true
}

type createUserRequest struct {
	Email    string                     `json:"email"`
	Password string                     `json:"password"`
	Role     string                     `json:"role"`
	Status   string                     `json:"status"`
	Profile  *updateUserProfileRequest  `json:"profile"`
	Settings *updateUserSettingsRequest `json:"settings"`
}

type updateUserAdminRequest struct {
	Email    string                     `json:"email"`
	Password string                     `json:"password"`
	Role     string                     `json:"role"`
	Status   string                     `json:"status"`
	Profile  *updateUserProfileRequest  `json:"profile"`
	Settings *updateUserSettingsRequest `json:"settings"`
}

func (h *UserHandler) List(c *gin.Context) {
	page := parseIntQueryUser(c, "page", 1)
	pageSize := parseIntQueryUser(c, "page_size", 20)
	search := c.Query("search")
	role := c.Query("role")
	status := c.Query("status")

	users, total, err := h.userSvc.List(c.Request.Context(), user.ListInput{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
		Role:     role,
		Status:   status,
	})
	if err != nil {
		respondInternal(c, "failed to list users")
		return
	}

	items := make([]gin.H, 0, len(users))
	for _, userItem := range users {
		displayName := ""
		if userItem.Profile != nil {
			if userItem.Profile.DisplayName != nil && *userItem.Profile.DisplayName != "" {
				displayName = *userItem.Profile.DisplayName
			} else if userItem.Profile.FullName != nil && *userItem.Profile.FullName != "" {
				displayName = *userItem.Profile.FullName
			}
		}
		if displayName == "" {
			displayName = userItem.Email
		}
		items = append(items, gin.H{
			"id":            userItem.ID,
			"email":         userItem.Email,
			"username":      displayName,
			"status":        userItem.Status,
			"role":          userItem.Role,
			"last_login_at": userItem.LastLoginAt,
			"created_at":    userItem.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     items,
		"page":      page,
		"page_size": pageSize,
		"total":     total,
	})
}

func (h *UserHandler) Create(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	var profile *user.UpdateProfile
	if req.Profile != nil {
		profile = &user.UpdateProfile{
			DisplayName: req.Profile.DisplayName,
			FullName:    req.Profile.FullName,
			Title:       req.Profile.Title,
			AvatarURL:   req.Profile.AvatarURL,
			Bio:         req.Profile.Bio,
			Phone:       req.Profile.Phone,
		}
	}

	var settings *user.UpdateSettings
	if req.Settings != nil {
		settings = &user.UpdateSettings{
			Theme:    req.Settings.Theme,
			Locale:   req.Settings.Locale,
			Timezone: req.Settings.Timezone,
		}
	}

	userModel, err := h.userSvc.Create(c.Request.Context(), user.CreateInput{
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
		Status:   req.Status,
		Profile:  profile,
		Settings: settings,
	})
	if err != nil {
		if err == user.ErrEmailTaken {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email taken"})
			return
		}
		respondInternal(c, "failed to create user")
		return
	}

	c.JSON(http.StatusCreated, userModel)
}

func (h *UserHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userModel, profile, settings, err := h.userSvc.Get(c.Request.Context(), id)
	if err != nil {
		if err == user.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		respondInternal(c, "failed to get user")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":            userModel.ID,
			"email":         userModel.Email,
			"status":        userModel.Status,
			"role":          userModel.Role,
			"last_login_at": userModel.LastLoginAt,
			"created_at":    userModel.CreatedAt,
		},
		"profile":  profile,
		"settings": settings,
	})
}

func (h *UserHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req updateUserAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	var profile *user.UpdateProfile
	if req.Profile != nil {
		profile = &user.UpdateProfile{
			DisplayName: req.Profile.DisplayName,
			FullName:    req.Profile.FullName,
			Title:       req.Profile.Title,
			AvatarURL:   req.Profile.AvatarURL,
			Bio:         req.Profile.Bio,
			Phone:       req.Profile.Phone,
		}
	}

	var settings *user.UpdateSettings
	if req.Settings != nil {
		settings = &user.UpdateSettings{
			Theme:    req.Settings.Theme,
			Locale:   req.Settings.Locale,
			Timezone: req.Settings.Timezone,
		}
	}

	if err := h.userSvc.Update(c.Request.Context(), id, req.Email, req.Password, req.Role, req.Status, profile, settings); err != nil {
		if err == user.ErrEmailTaken {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email taken"})
			return
		}
		if err == user.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		respondInternal(c, "failed to update user")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *UserHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.userSvc.Delete(c.Request.Context(), id); err != nil {
		respondInternal(c, "failed to delete user")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func parseIntQueryUser(c *gin.Context, key string, fallback int) int {
	value := c.Query(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
