package handlers

import (
	"errors"
	"net/http"

	"deepspace/internal/service/email"
	"deepspace/internal/service/passwordreset"

	"github.com/gin-gonic/gin"
)

type PasswordResetHandler struct {
	svc *passwordreset.Service
}

func NewPasswordResetHandler(svc *passwordreset.Service) *PasswordResetHandler {
	return &PasswordResetHandler{svc: svc}
}

type passwordResetRequest struct {
	Email string `json:"email"`
}

type passwordResetConfirmRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

// RequestPasswordReset godoc
// @Summary 申请重置密码
// @Description 发送重置密码邮件
// @Tags 认证
// @Accept json
// @Produce json
// @Param data body passwordResetRequest true "重置请求"
// @Success 200 {object} map[string]interface{} "请求成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 503 {object} map[string]interface{} "服务不可用"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /auth/password-reset/request [post]
func (h *PasswordResetHandler) RequestPasswordReset(c *gin.Context) {
	if h == nil || h.svc == nil {
		respondInternal(c, "重置服务未配置")
		return
	}

	var req passwordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数不正确"})
		return
	}

	if err := h.svc.RequestReset(c.Request.Context(), req.Email); err != nil {
		switch {
		case errors.Is(err, passwordreset.ErrInvalidEmail):
			c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱不正确"})
			return
		case errors.Is(err, passwordreset.ErrRedisDisabled):
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "重置服务不可用"})
			return
		case errors.Is(err, passwordreset.ErrMissingBaseURL):
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "重置服务未配置"})
			return
		case errors.Is(err, email.ErrEmailDisabled):
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "邮件服务未启用"})
			return
		default:
			respondInternal(c, "发送重置邮件失败")
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// ConfirmPasswordReset godoc
// @Summary 确认重置密码
// @Description 校验重置令牌并更新密码
// @Tags 认证
// @Accept json
// @Produce json
// @Param data body passwordResetConfirmRequest true "重置确认"
// @Success 200 {object} map[string]interface{} "重置成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 503 {object} map[string]interface{} "服务不可用"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /auth/password-reset/confirm [post]
func (h *PasswordResetHandler) ConfirmPasswordReset(c *gin.Context) {
	if h == nil || h.svc == nil {
		respondInternal(c, "重置服务未配置")
		return
	}

	var req passwordResetConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数不正确"})
		return
	}

	if err := h.svc.ConfirmReset(c.Request.Context(), req.Token, req.NewPassword); err != nil {
		switch {
		case errors.Is(err, passwordreset.ErrInvalidToken):
			c.JSON(http.StatusBadRequest, gin.H{"error": "重置令牌无效或已过期"})
			return
		case errors.Is(err, passwordreset.ErrInvalidPassword):
			c.JSON(http.StatusBadRequest, gin.H{"error": "密码至少 8 位字符"})
			return
		case errors.Is(err, passwordreset.ErrRedisDisabled):
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "重置服务不可用"})
			return
		default:
			respondInternal(c, "重置密码失败")
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
