package handlers

import (
	"net/http"

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
			"id":     userModel.ID,
			"email":  userModel.Email,
			"status": userModel.Status,
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
			"id":     userModel.ID,
			"email":  userModel.Email,
			"status": userModel.Status,
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
