package repo

import (
	"context"
	"errors"

	"deepspace/internal/model"

	"gorm.io/gorm"
)

type RiskPolicyRepo struct {
	db *gorm.DB
}

func NewRiskPolicyRepo(db *gorm.DB) *RiskPolicyRepo {
	return &RiskPolicyRepo{db: db}
}

type RiskPolicyFilter struct {
	Scope     string
	UserID    *int64
	ProjectID *int64
	Status    string
	Limit     int
	Offset    int
}

func (r *RiskPolicyRepo) Create(ctx context.Context, item *model.RiskPolicy) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *RiskPolicyRepo) Update(ctx context.Context, id int64, updates map[string]any) (*model.RiskPolicy, error) {
	if len(updates) == 0 {
		return r.GetByID(ctx, id)
	}
	if err := r.db.WithContext(ctx).
		Model(&model.RiskPolicy{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *RiskPolicyRepo) Delete(ctx context.Context, id int64) (bool, error) {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.RiskPolicy{})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *RiskPolicyRepo) GetByID(ctx context.Context, id int64) (*model.RiskPolicy, error) {
	var item model.RiskPolicy
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *RiskPolicyRepo) List(ctx context.Context, filter RiskPolicyFilter) ([]model.RiskPolicy, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.RiskPolicy{})
	if filter.Scope != "" {
		query = query.Where("scope = ?", filter.Scope)
	}
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.ProjectID != nil {
		query = query.Where("project_id = ?", *filter.ProjectID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []model.RiskPolicy
	if err := query.Order("priority ASC, id DESC").Limit(filter.Limit).Offset(filter.Offset).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
