package user

import (
	"context"
	"errors"
	"strings"

	"deepspace/internal/model"
	"deepspace/internal/repo"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrInvalidSetting = errors.New("invalid settings")
	ErrEmailTaken     = errors.New("email already taken")
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

type ListInput struct {
	Page     int
	PageSize int
	Search   string
	Role     string
	Status   string
}

type CreateInput struct {
	Email    string
	Password string
	Role     string
	Status   string
	Profile  *UpdateProfile
	Settings *UpdateSettings
}

func (s *Service) List(ctx context.Context, input ListInput) ([]model.User, int64, error) {
	return s.users.List(ctx, input.Page, input.PageSize, input.Search, input.Role, input.Status)
}

func (s *Service) Get(ctx context.Context, id int64) (*model.User, *model.UserProfile, *model.UserSettings, error) {
	return s.GetMe(ctx, id)
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*model.User, error) {
	existing, err := s.users.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:        input.Email,
		PasswordHash: string(hash),
		Status:       input.Status,
		Role:         input.Role,
	}
	if user.Status == "" {
		user.Status = "active"
	}
	if user.Role == "" {
		user.Role = "user"
	}

	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}

	// Initialize profile
	profile := &model.UserProfile{UserID: user.ID}
	if input.Profile != nil {
		profile.DisplayName = input.Profile.DisplayName
		profile.FullName = input.Profile.FullName
		profile.Title = input.Profile.Title
		profile.AvatarURL = input.Profile.AvatarURL
		profile.Bio = input.Profile.Bio
		profile.Phone = input.Profile.Phone
	}
	if err := s.profiles.Upsert(ctx, profile); err != nil {
		return nil, err
	}

	// Initialize settings
	settings := &model.UserSettings{
		UserID:   user.ID,
		Theme:    DefaultTheme,
		Locale:   DefaultLocale,
		Timezone: DefaultTimezone,
	}
	if input.Settings != nil {
		if input.Settings.Theme != nil {
			settings.Theme = *input.Settings.Theme
		}
		if input.Settings.Locale != nil {
			settings.Locale = *input.Settings.Locale
		}
		if input.Settings.Timezone != nil {
			settings.Timezone = *input.Settings.Timezone
		}
	}
	if err := s.settings.Upsert(ctx, settings); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) Update(ctx context.Context, id int64, email string, password string, role string, status string, profile *UpdateProfile, settings *UpdateSettings) error {
	user, _, _, err := s.GetMe(ctx, id)
	if err != nil {
		return err
	}

	// Update User fields
	if email != "" && email != user.Email {
		existing, err := s.users.GetByEmail(ctx, email)
		if err != nil {
			return err
		}
		if existing != nil && existing.ID != id {
			return ErrEmailTaken
		}
		user.Email = email
	}

	if password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.PasswordHash = string(hash)
	}

	if role != "" {
		user.Role = role
	}

	if status != "" {
		user.Status = status
	}

	if err := s.users.Update(ctx, user); err != nil {
		return err
	}

	// Reuse UpdateMe logic for Profile/Settings by calling it (UpdateMe updates profile/settings in DB)
	// But UpdateMe retrieves user internally. We can just call UpdateMe logic or extract it.
	// Since we already have the profile/settings structs to update, we can just use UpdateMe logic.
	// However, UpdateMe takes structs with pointers.

	_, _, _, err = s.UpdateMe(ctx, id, profile, settings)
	return err
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.users.Delete(ctx, id)
}
