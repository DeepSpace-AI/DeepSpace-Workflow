package auth

import (
	"context"
	"errors"
	"strings"

	"deepspace/internal/repo"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailTaken         = errors.New("email already registered")
)

type UserAuthService struct {
	users *repo.UserRepo
	orgs  *repo.OrgRepo
	jwt   *JWTManager
}

func NewUserAuthService(users *repo.UserRepo, orgs *repo.OrgRepo, jwt *JWTManager) *UserAuthService {
	return &UserAuthService{users: users, orgs: orgs, jwt: jwt}
}

type AuthResult struct {
	UserID int64
	OrgID  int64
	Token  string
}

func (s *UserAuthService) Register(ctx context.Context, email, password, orgName string) (*AuthResult, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || password == "" {
		return nil, ErrInvalidCredentials
	}
	if orgName == "" {
		orgName = "Default Org"
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

	user, err := s.users.Create(ctx, email, string(hash))
	if err != nil {
		return nil, err
	}

	org, err := s.orgs.Create(ctx, orgName, user.ID)
	if err != nil {
		return nil, err
	}
	if err := s.orgs.AddMember(ctx, org.ID, user.ID, "owner"); err != nil {
		return nil, err
	}

	token, err := s.jwt.Sign(user.ID, org.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResult{UserID: user.ID, OrgID: org.ID, Token: token}, nil
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

	orgID := int64(0)
	org, err := s.orgs.GetByOwner(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	if org == nil {
		return nil, ErrInvalidCredentials
	}
	orgID = org.ID

	token, err := s.jwt.Sign(user.ID, orgID)
	if err != nil {
		return nil, err
	}

	return &AuthResult{UserID: user.ID, OrgID: orgID, Token: token}, nil
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
