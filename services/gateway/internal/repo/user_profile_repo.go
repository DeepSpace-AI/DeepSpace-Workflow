package repo

import (
	"context"
	"errors"

	"deepspace/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserProfileRepo struct {
	db *gorm.DB
}

func NewUserProfileRepo(db *gorm.DB) *UserProfileRepo {
	return &UserProfileRepo{db: db}
}

func (r *UserProfileRepo) GetByUserID(ctx context.Context, userID int64) (*model.UserProfile, error) {
	var profile model.UserProfile
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &profile, nil
}

func (r *UserProfileRepo) Upsert(ctx context.Context, profile *model.UserProfile) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"display_name", "full_name", "title", "avatar_url", "bio", "phone", "updated_at"}),
	}).Create(profile).Error
}
