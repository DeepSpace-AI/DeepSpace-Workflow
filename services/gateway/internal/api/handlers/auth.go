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

func (h *AuthHandler) Logout(c *gin.Context) {
	clearAuthCookie(c, h.jwt.CookieName, h.jwt)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

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
