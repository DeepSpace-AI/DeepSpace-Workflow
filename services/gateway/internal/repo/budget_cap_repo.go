package repo

import (
	"context"
	"errors"

	"deepspace/internal/model"

	"gorm.io/gorm"
)

type BudgetCapRepo struct {
	db *gorm.DB
}

func NewBudgetCapRepo(db *gorm.DB) *BudgetCapRepo {
	return &BudgetCapRepo{db: db}
}

type BudgetCapFilter struct {
	PolicyID *int64
	Status   string
	Limit    int
	Offset   int
}

func (r *BudgetCapRepo) Create(ctx context.Context, item *model.BudgetCap) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *BudgetCapRepo) Update(ctx context.Context, id int64, updates map[string]any) (*model.BudgetCap, error) {
	if len(updates) == 0 {
		return r.GetByID(ctx, id)
	}
	if err := r.db.WithContext(ctx).
		Model(&model.BudgetCap{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *BudgetCapRepo) Delete(ctx context.Context, id int64) (bool, error) {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.BudgetCap{})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *BudgetCapRepo) GetByID(ctx context.Context, id int64) (*model.BudgetCap, error) {
	var item model.BudgetCap
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *BudgetCapRepo) List(ctx context.Context, filter BudgetCapFilter) ([]model.BudgetCap, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.BudgetCap{})
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

	var items []model.BudgetCap
	if err := query.Order("id DESC").Limit(filter.Limit).Offset(filter.Offset).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
