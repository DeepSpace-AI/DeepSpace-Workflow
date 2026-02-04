package repo

import (
	"context"

	"deepspace/internal/model"

	"gorm.io/gorm"
)

type UsageRepo struct {
	db *gorm.DB
}

func NewUsageRepo(db *gorm.DB) *UsageRepo {
	return &UsageRepo{db: db}
}

func (r *UsageRepo) Create(ctx context.Context, record model.UsageRecord) error {
	return r.db.WithContext(ctx).Create(&record).Error
}
