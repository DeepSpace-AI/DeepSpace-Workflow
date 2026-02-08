package repo

import (
	"context"
	"errors"

	"deepspace/internal/model"

	"gorm.io/gorm"
)

type PlanModelPriceRepo struct {
	db *gorm.DB
}

func NewPlanModelPriceRepo(db *gorm.DB) *PlanModelPriceRepo {
	return &PlanModelPriceRepo{db: db}
}

func (r *PlanModelPriceRepo) Create(ctx context.Context, item *model.PlanModelPrice) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *PlanModelPriceRepo) Update(ctx context.Context, id int64, updates map[string]any) (*model.PlanModelPrice, error) {
	if len(updates) == 0 {
		return r.GetByID(ctx, id)
	}
	if err := r.db.WithContext(ctx).
		Model(&model.PlanModelPrice{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *PlanModelPriceRepo) GetByID(ctx context.Context, id int64) (*model.PlanModelPrice, error) {
	var item model.PlanModelPrice
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

func (r *PlanModelPriceRepo) GetByPlanModel(ctx context.Context, planID int64, modelID string) (*model.PlanModelPrice, error) {
	var item model.PlanModelPrice
	err := r.db.WithContext(ctx).
		Where("plan_id = ? AND model_id = ?", planID, modelID).
		First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *PlanModelPriceRepo) ListByPlanID(ctx context.Context, planID int64) ([]model.PlanModelPrice, error) {
	var items []model.PlanModelPrice
	if err := r.db.WithContext(ctx).
		Where("plan_id = ?", planID).
		Order("created_at DESC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
