package repo

import (
	"context"
	"errors"

	"deepspace/internal/model"

	"gorm.io/gorm"
)

type IPRuleRepo struct {
	db *gorm.DB
}

func NewIPRuleRepo(db *gorm.DB) *IPRuleRepo {
	return &IPRuleRepo{db: db}
}

type IPRuleFilter struct {
	PolicyID *int64
	Type     string
	Status   string
	Limit    int
	Offset   int
}

func (r *IPRuleRepo) Create(ctx context.Context, item *model.IPRule) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *IPRuleRepo) Update(ctx context.Context, id int64, updates map[string]any) (*model.IPRule, error) {
	if len(updates) == 0 {
		return r.GetByID(ctx, id)
	}
	if err := r.db.WithContext(ctx).
		Model(&model.IPRule{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *IPRuleRepo) Delete(ctx context.Context, id int64) (bool, error) {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.IPRule{})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *IPRuleRepo) GetByID(ctx context.Context, id int64) (*model.IPRule, error) {
	var item model.IPRule
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *IPRuleRepo) List(ctx context.Context, filter IPRuleFilter) ([]model.IPRule, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.IPRule{})
	if filter.PolicyID != nil {
		query = query.Where("policy_id = ?", *filter.PolicyID)
	}
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []model.IPRule
	if err := query.Order("id DESC").Limit(filter.Limit).Offset(filter.Offset).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
