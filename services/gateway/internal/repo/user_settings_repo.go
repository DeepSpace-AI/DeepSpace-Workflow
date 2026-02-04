package repo

import (
	"context"
	"errors"

	"deepspace/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserSettingsRepo struct {
	db *gorm.DB
}

func NewUserSettingsRepo(db *gorm.DB) *UserSettingsRepo {
	return &UserSettingsRepo{db: db}
}

func (r *UserSettingsRepo) GetByUserID(ctx context.Context, userID int64) (*model.UserSettings, error) {
	var settings model.UserSettings
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&settings).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &settings, nil
}

func (r *UserSettingsRepo) Upsert(ctx context.Context, settings *model.UserSettings) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"theme", "locale", "timezone", "updated_at"}),
	}).Create(settings).Error
}
