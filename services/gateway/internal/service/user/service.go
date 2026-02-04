package user

import (
	"context"
	"errors"
	"strings"

	"deepspace/internal/model"
	"deepspace/internal/repo"
)

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrInvalidSetting = errors.New("invalid settings")
)

const (
	DefaultTheme    = "system"
	DefaultLocale   = "zh-CN"
	DefaultTimezone = "Asia/Shanghai"
)

type Service struct {
	users    *repo.UserRepo
	profiles *repo.UserProfileRepo
	settings *repo.UserSettingsRepo
}

func New(users *repo.UserRepo, profiles *repo.UserProfileRepo, settings *repo.UserSettingsRepo) *Service {
	return &Service{users: users, profiles: profiles, settings: settings}
}

type UpdateProfile struct {
	DisplayName *string
	FullName    *string
	Title       *string
	AvatarURL   *string
	Bio         *string
	Phone       *string
}

type UpdateSettings struct {
	Theme    *string
	Locale   *string
	Timezone *string
}

func (s *Service) GetMe(ctx context.Context, userID int64) (*model.User, *model.UserProfile, *model.UserSettings, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, nil, nil, err
	}
	if user == nil {
		return nil, nil, nil, ErrUserNotFound
	}

	profile, err := s.profiles.GetByUserID(ctx, userID)
	if err != nil {
		return nil, nil, nil, err
	}
	if profile == nil {
		profile = &model.UserProfile{UserID: userID}
		if err := s.profiles.Upsert(ctx, profile); err != nil {
			return nil, nil, nil, err
		}
	}

	settings, err := s.settings.GetByUserID(ctx, userID)
	if err != nil {
		return nil, nil, nil, err
	}
	if settings == nil {
		settings = &model.UserSettings{
			UserID:   userID,
			Theme:    DefaultTheme,
			Locale:   DefaultLocale,
			Timezone: DefaultTimezone,
		}
		if err := s.settings.Upsert(ctx, settings); err != nil {
			return nil, nil, nil, err
		}
	}

	return user, profile, settings, nil
}

func (s *Service) UpdateMe(ctx context.Context, userID int64, profileUpdate *UpdateProfile, settingsUpdate *UpdateSettings) (*model.User, *model.UserProfile, *model.UserSettings, error) {
	user, profile, settings, err := s.GetMe(ctx, userID)
	if err != nil {
		return nil, nil, nil, err
	}

	if profileUpdate != nil {
		if profileUpdate.DisplayName != nil {
			profile.DisplayName = profileUpdate.DisplayName
		}
		if profileUpdate.FullName != nil {
			profile.FullName = profileUpdate.FullName
		}
		if profileUpdate.Title != nil {
			profile.Title = profileUpdate.Title
		}
		if profileUpdate.AvatarURL != nil {
			profile.AvatarURL = profileUpdate.AvatarURL
		}
		if profileUpdate.Bio != nil {
			profile.Bio = profileUpdate.Bio
		}
		if profileUpdate.Phone != nil {
			profile.Phone = profileUpdate.Phone
		}
		if err := s.profiles.Upsert(ctx, profile); err != nil {
			return nil, nil, nil, err
		}
	}

	if settingsUpdate != nil {
		if settingsUpdate.Theme != nil {
			value := strings.TrimSpace(*settingsUpdate.Theme)
			if !isValidTheme(value) {
				return nil, nil, nil, ErrInvalidSetting
			}
			settings.Theme = value
		}
		if settingsUpdate.Locale != nil {
			value := strings.TrimSpace(*settingsUpdate.Locale)
			if !isValidLocale(value) {
				return nil, nil, nil, ErrInvalidSetting
			}
			settings.Locale = value
		}
		if settingsUpdate.Timezone != nil {
			value := strings.TrimSpace(*settingsUpdate.Timezone)
			if value == "" {
				return nil, nil, nil, ErrInvalidSetting
			}
			settings.Timezone = value
		}
		if err := s.settings.Upsert(ctx, settings); err != nil {
			return nil, nil, nil, err
		}
	}

	return user, profile, settings, nil
}

func isValidTheme(value string) bool {
	switch value {
	case "light", "dark", "system":
		return true
	default:
		return false
	}
}

func isValidLocale(value string) bool {
	switch value {
	case "zh-CN", "en-US":
		return true
	default:
		return false
	}
}
