package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"deepspace/internal/model"
	"deepspace/internal/repo"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailTaken         = errors.New("email already registered")
)

type UserAuthService struct {
	users *repo.UserRepo
	jwt   *JWTManager
}

func NewUserAuthService(users *repo.UserRepo, jwt *JWTManager) *UserAuthService {
	return &UserAuthService{users: users, jwt: jwt}
}

type AuthResult struct {
	UserID int64
	Token  string
}

func (s *UserAuthService) Register(ctx context.Context, email, password string) (*AuthResult, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || password == "" {
		return nil, ErrInvalidCredentials
	}

	existing, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:        email,
		PasswordHash: string(hash),
		Status:       "active",
		Role:         "user",
	}
	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}

	token, err := s.jwt.Sign(user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResult{UserID: user.ID, Token: token}, nil
}

func (s *UserAuthService) Login(ctx context.Context, email, password string) (*AuthResult, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Update last login time
	if err := s.users.UpdateLastLogin(ctx, user.ID, time.Now()); err != nil {
		// Log error but don't fail login? Or fail?
		// For now, let's just proceed or log. Since we don't have a logger here, we can ignore or return error.
		// Returning error might be strict but safe.
		// Let's ignore it for now to avoid login failure due to non-critical update.
		// actually, let's just ignore the error for now as it is not critical
	}

	token, err := s.jwt.Sign(user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResult{UserID: user.ID, Token: token}, nil
}

func (s *UserAuthService) ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error {
	if strings.TrimSpace(oldPassword) == "" || strings.TrimSpace(newPassword) == "" {
		return ErrInvalidCredentials
	}

	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return ErrInvalidCredentials
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.users.UpdatePassword(ctx, userID, string(hash))
}
