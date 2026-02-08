package repo

import (
	"context"
	"errors"

	"deepspace/internal/model"

	"gorm.io/gorm"
)

type PlanRepo struct {
	db *gorm.DB
}

func NewPlanRepo(db *gorm.DB) *PlanRepo {
	return &PlanRepo{db: db}
}

func (r *PlanRepo) Create(ctx context.Context, plan *model.Plan) error {
	return r.db.WithContext(ctx).Create(plan).Error
}

func (r *PlanRepo) Update(ctx context.Context, id int64, updates map[string]any) (*model.Plan, error) {
	if len(updates) == 0 {
		return r.GetByID(ctx, id)
	}
	if err := r.db.WithContext(ctx).
		Model(&model.Plan{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *PlanRepo) GetByID(ctx context.Context, id int64) (*model.Plan, error) {
	var plan model.Plan
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&plan).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &plan, nil
}

func (r *PlanRepo) List(ctx context.Context) ([]model.Plan, error) {
	var plans []model.Plan
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Find(&plans).Error; err != nil {
		return nil, err
	}
	return plans, nil
}
