package repo

import (
	"context"
	"errors"

	"deepspace/internal/model"

	"gorm.io/gorm"
)

type RateLimitRepo struct {
	db *gorm.DB
}

func NewRateLimitRepo(db *gorm.DB) *RateLimitRepo {
	return &RateLimitRepo{db: db}
}

type RateLimitFilter struct {
	PolicyID *int64
	Status   string
	Limit    int
	Offset   int
}

func (r *RateLimitRepo) Create(ctx context.Context, item *model.RateLimit) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *RateLimitRepo) Update(ctx context.Context, id int64, updates map[string]any) (*model.RateLimit, error) {
	if len(updates) == 0 {
		return r.GetByID(ctx, id)
	}
	if err := r.db.WithContext(ctx).
		Model(&model.RateLimit{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *RateLimitRepo) Delete(ctx context.Context, id int64) (bool, error) {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.RateLimit{})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *RateLimitRepo) GetByID(ctx context.Context, id int64) (*model.RateLimit, error) {
	var item model.RateLimit
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *RateLimitRepo) List(ctx context.Context, filter RateLimitFilter) ([]model.RateLimit, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.RateLimit{})
	if filter.PolicyID != nil {
		query = query.Where("policy_id = ?", *filter.PolicyID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []model.RateLimit
	if err := query.Order("id DESC").Limit(filter.Limit).Offset(filter.Offset).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
