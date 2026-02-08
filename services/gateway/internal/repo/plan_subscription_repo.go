package repo

import (
	"context"
	"errors"
	"time"

	"deepspace/internal/model"

	"gorm.io/gorm"
)

type PlanSubscriptionRepo struct {
	db *gorm.DB
}

func NewPlanSubscriptionRepo(db *gorm.DB) *PlanSubscriptionRepo {
	return &PlanSubscriptionRepo{db: db}
}

func (r *PlanSubscriptionRepo) Create(ctx context.Context, item *model.PlanSubscription) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *PlanSubscriptionRepo) Update(ctx context.Context, id int64, updates map[string]any) (*model.PlanSubscription, error) {
	if len(updates) == 0 {
		return r.GetByID(ctx, id)
	}
	if err := r.db.WithContext(ctx).
		Model(&model.PlanSubscription{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *PlanSubscriptionRepo) GetByID(ctx context.Context, id int64) (*model.PlanSubscription, error) {
	var item model.PlanSubscription
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *PlanSubscriptionRepo) ListByOrg(ctx context.Context, orgID int64) ([]model.PlanSubscription, error) {
	var items []model.PlanSubscription
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", orgID).
		Order("created_at DESC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *PlanSubscriptionRepo) GetActiveByOrg(ctx context.Context, orgID int64, now time.Time) (*model.PlanSubscription, error) {
	var item model.PlanSubscription
	err := r.db.WithContext(ctx).
		Where("user_id = ?", orgID).
		Where("status = ?", "active").
		Where("start_at <= ?", now).
		Where("end_at IS NULL OR end_at > ?", now).
		Order("start_at DESC").
		First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}
