package passwordreset

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/url"
	"strconv"
	"strings"
	"time"

	"deepspace/internal/config"
	"deepspace/internal/repo"
	"deepspace/internal/service/email"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

const (
	defaultTokenBytes = 32
	defaultTokenTTL   = 30 * time.Minute
	resetPath         = "/reset-password"
)

var (
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidToken    = errors.New("invalid token")
	ErrInvalidPassword = errors.New("invalid password")
	ErrRedisDisabled   = errors.New("redis disabled")
	ErrMissingBaseURL  = errors.New("missing web base url")
)

type Service struct {
	cfg      *config.Config
	users    *repo.UserRepo
	profiles *repo.UserProfileRepo
	emailSvc *email.Service
	redis    *redis.Client
	tokenTTL time.Duration
}

func New(cfg *config.Config, users *repo.UserRepo, profiles *repo.UserProfileRepo, emailSvc *email.Service) (*Service, error) {
	if cfg == nil || users == nil || profiles == nil || emailSvc == nil {
		return nil, errors.New("missing dependency")
	}

	svc := &Service{
		cfg:      cfg,
		users:    users,
		profiles: profiles,
		emailSvc: emailSvc,
		tokenTTL: defaultTokenTTL,
	}

	if strings.TrimSpace(cfg.RedisURL) != "" {
		opt, err := redis.ParseURL(cfg.RedisURL)
		if err != nil {
			return nil, err
		}
		svc.redis = redis.NewClient(opt)
	}

	return svc, nil
}

func (s *Service) RequestReset(ctx context.Context, emailAddr string) error {
	emailAddr = strings.TrimSpace(strings.ToLower(emailAddr))
	if emailAddr == "" {
		return ErrInvalidEmail
	}
	if s.redis == nil {
		return ErrRedisDisabled
	}

	user, err := s.users.GetByEmail(ctx, emailAddr)
	if err != nil {
		return err
	}
	if user == nil {
		return nil
	}

	token, tokenHash, err := generateToken()
	if err != nil {
		return err
	}

	key := s.redisKey(tokenHash)
	if err := s.redis.Set(ctx, key, strconv.FormatInt(user.ID, 10), s.tokenTTL).Err(); err != nil {
		return err
	}

	resetURL, err := s.buildResetURL(token)
	if err != nil {
		return err
	}

	username := s.resolveUsername(ctx, user.ID, emailAddr)
	input := email.EmailInput{
		Type:    email.EmailTypeResetPassword,
		To:      []string{emailAddr},
		Subject: "重置密码",
		TemplateData: map[string]any{
			"username": username,
			"date":     time.Now().Format("2006-01-02 15:04:05"),
			"address":  resetURL,
		},
	}

	return s.emailSvc.Send(ctx, input)
}

func (s *Service) ConfirmReset(ctx context.Context, token, newPassword string) error {
	token = strings.TrimSpace(token)
	newPassword = strings.TrimSpace(newPassword)
	if token == "" {
		return ErrInvalidToken
	}
	if len(newPassword) < 8 {
		return ErrInvalidPassword
	}
	if s.redis == nil {
		return ErrRedisDisabled
	}

	key := s.redisKey(hashToken(token))
	value, err := s.redis.GetDel(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrInvalidToken
		}
		return err
	}

	userID, err := strconv.ParseInt(value, 10, 64)
	if err != nil || userID <= 0 {
		return ErrInvalidToken
	}

	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrInvalidToken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.users.UpdatePassword(ctx, userID, string(hash))
}

func (s *Service) resolveUsername(ctx context.Context, userID int64, emailAddr string) string {
	profile, err := s.profiles.GetByUserID(ctx, userID)
	if err == nil && profile != nil {
		if profile.DisplayName != nil {
			value := strings.TrimSpace(*profile.DisplayName)
			if value != "" {
				return value
			}
		}
		if profile.FullName != nil {
			value := strings.TrimSpace(*profile.FullName)
			if value != "" {
				return value
			}
		}
	}
	return emailAddr
}

func (s *Service) buildResetURL(token string) (string, error) {
	base := strings.TrimSpace(s.cfg.WebBaseURL)
	if base == "" {
		return "", ErrMissingBaseURL
	}
	base = strings.TrimRight(base, "/")
	return base + resetPath + "?token=" + url.QueryEscape(token), nil
}

func (s *Service) redisKey(hash string) string {
	return "password_reset:" + hash
}

func generateToken() (string, string, error) {
	buf := make([]byte, defaultTokenBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", "", err
	}
	token := hex.EncodeToString(buf)
	return token, hashToken(token), nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
