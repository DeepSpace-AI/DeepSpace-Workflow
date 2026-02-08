package handlers

import (
	"net/http"
	"time"

	"deepspace/internal/service/auth"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc *auth.UserAuthService
	jwt *auth.JWTManager
}

func NewAuthHandler(svc *auth.UserAuthService, jwt *auth.JWTManager) *AuthHandler {
	return &AuthHandler{svc: svc, jwt: jwt}
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register godoc
// @Summary 用户注册
// @Description 使用邮箱密码注册并设置登录 Cookie
// @Tags 认证
// @Accept json
// @Produce json
// @Param data body registerRequest true "注册信息"
// @Success 201 {object} map[string]interface{} "创建成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 409 {object} map[string]interface{} "邮箱已注册"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	result, err := h.svc.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch err {
		case auth.ErrEmailTaken:
			c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
			return
		case auth.ErrInvalidCredentials:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials"})
			return
		default:
			respondInternal(c, "registration failed")
			return
		}
	}

	setAuthCookie(c, h.jwt.CookieName, result.Token, h.jwt.ExpiresIn, h.jwt)
	c.JSON(http.StatusCreated, gin.H{"user_id": result.UserID})
}

// Login godoc
// @Summary 用户登录
// @Description 使用邮箱密码登录并设置登录 Cookie
// @Tags 认证
// @Accept json
// @Produce json
// @Param data body loginRequest true "登录信息"
// @Success 200 {object} map[string]interface{} "登录成功"
// @Failure 400 {object} map[string]interface{} "请求错误"
// @Failure 401 {object} map[string]interface{} "账号或密码错误"
// @Failure 500 {object} map[string]interface{} "服务内部错误"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	result, err := h.svc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch err {
		case auth.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		default:
			respondInternal(c, "login failed")
			return
		}
	}

	setAuthCookie(c, h.jwt.CookieName, result.Token, h.jwt.ExpiresIn, h.jwt)
	c.JSON(http.StatusOK, gin.H{"user_id": result.UserID})
}

// Logout godoc
// @Summary 用户退出登录
// @Description 清理登录 Cookie
// @Tags 认证
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "退出成功"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	clearAuthCookie(c, h.jwt.CookieName, h.jwt)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Me godoc
// @Summary 获取当前用户
// @Description 获取当前登录用户的 user_id
// @Tags 认证
// @Accept json
// @Produce json
// @Security bearerAuth
// @Security cookieAuth
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 401 {object} map[string]interface{} "未登录"
// @Router /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
	})
}

func setAuthCookie(c *gin.Context, name, token string, expiresIn time.Duration, jwt *auth.JWTManager) {
	c.SetCookie(name, token, int(expiresIn.Seconds()), "/", "", jwtSecure(jwt), true)
}

func clearAuthCookie(c *gin.Context, name string, jwt *auth.JWTManager) {
	c.SetCookie(name, "", -1, "/", "", jwtSecure(jwt), true)
}

func jwtSecure(jwt *auth.JWTManager) bool {
	if jwt == nil {
		return false
	}
	return jwt.CookieSecure
}
