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

// GetMe godoc
// @Summary 获取当前用户信息
// @Description 返回当前登录用户的基础信息、资料与设置
// @Tags 用户
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 404 {object} map[string]interface{} "用户不存在"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /users/me [get]
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

// UpdateMe godoc
// @Summary 更新当前用户信息
// @Description 更新当前登录用户的资料与设置
// @Tags 用户
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param data body updateUserRequest true "用户资料与设置"
// @Success 200 {object} map[string]interface{} "更新成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 404 {object} map[string]interface{} "用户不存在"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /users/me [patch]
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

// ChangePassword godoc
// @Summary 修改当前用户密码
// @Description 校验旧密码后更新新密码
// @Tags 用户
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param data body changePasswordRequest true "密码信息"
// @Success 200 {object} map[string]interface{} "修改成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "账号或密码错误"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /users/me/password [post]
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

// List godoc
// @Summary 管理员：用户列表
// @Description 需要管理员权限
// @Tags 管理-用户
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param search query string false "搜索关键字"
// @Param role query string false "角色"
// @Param status query string false "状态"
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 403 {object} map[string]interface{} "无权限"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /admin/users [get]
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

// Create godoc
// @Summary 管理员：创建用户
// @Description 需要管理员权限
// @Tags 管理-用户
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param data body createUserRequest true "用户信息"
// @Success 201 {object} map[string]interface{} "创建成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 403 {object} map[string]interface{} "无权限"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /admin/users [post]
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

// Get godoc
// @Summary 管理员：获取用户
// @Description 需要管理员权限
// @Tags 管理-用户
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param id path int true "用户ID"
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 403 {object} map[string]interface{} "无权限"
// @Failure 404 {object} map[string]interface{} "用户不存在"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /admin/users/{id} [get]
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

// Update godoc
// @Summary 管理员：更新用户
// @Description 需要管理员权限
// @Tags 管理-用户
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param id path int true "用户ID"
// @Param data body updateUserAdminRequest true "用户信息"
// @Success 200 {object} map[string]interface{} "更新成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 403 {object} map[string]interface{} "无权限"
// @Failure 404 {object} map[string]interface{} "用户不存在"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /admin/users/{id} [patch]
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

// Delete godoc
// @Summary 管理员：删除用户
// @Description 需要管理员权限
// @Tags 管理-用户
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Param id path int true "用户ID"
// @Success 200 {object} map[string]interface{} "删除成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Failure 403 {object} map[string]interface{} "无权限"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /admin/users/{id} [delete]
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
