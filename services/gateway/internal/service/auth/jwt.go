package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	Secret       []byte
	Issuer       string
	ExpiresIn    time.Duration
	CookieName   string
	CookieSecure bool
}

type Claims struct {
	UserID int64 `json:"user_id"`
	OrgID  int64 `json:"org_id"`
	jwt.RegisteredClaims
}

func (m *JWTManager) Sign(userID, orgID int64) (string, error) {
	claims := Claims{
		UserID: userID,
		OrgID:  orgID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.Issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.ExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.Secret)
}

func (m *JWTManager) Verify(tokenString string) (*Claims, error) {
	if tokenString == "" {
		return nil, errors.New("empty token")
	}
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return m.Secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
