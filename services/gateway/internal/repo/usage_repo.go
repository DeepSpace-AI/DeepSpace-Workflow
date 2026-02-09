package repo

import (
	"context"
	"time"

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

type UsageListFilter struct {
	UserID int64
	Start  *time.Time
	End    *time.Time
	Limit  int
	Offset int
}

func (r *UsageRepo) List(ctx context.Context, filter UsageListFilter) ([]model.UsageRecord, error) {
	query := r.db.WithContext(ctx).
		Where("user_id = ?", filter.UserID)
	if filter.Start != nil {
		query = query.Where("created_at >= ?", *filter.Start)
	}
	if filter.End != nil {
		query = query.Where("created_at < ?", *filter.End)
	}

	var records []model.UsageRecord
	if err := query.
		Order("created_at DESC").
		Limit(filter.Limit).
		Offset(filter.Offset).
		Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func (r *UsageRepo) Count(ctx context.Context, userID int64, start, end *time.Time) (int64, error) {
	query := r.db.WithContext(ctx).
		Model(&model.UsageRecord{}).
		Where("user_id = ?", userID)
	if start != nil {
		query = query.Where("created_at >= ?", *start)
	}
	if end != nil {
		query = query.Where("created_at < ?", *end)
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *UsageRepo) SumCost(ctx context.Context, userID int64, start, end *time.Time) (float64, error) {
	query := r.db.WithContext(ctx).
		Model(&model.UsageRecord{}).
		Where("user_id = ?", userID)
	if start != nil {
		query = query.Where("created_at >= ?", *start)
	}
	if end != nil {
		query = query.Where("created_at < ?", *end)
	}
	var total float64
	if err := query.Select("COALESCE(SUM(cost), 0)").Scan(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

type AdminUsageListFilter struct {
	UserID *int64
	Start  *time.Time
	End    *time.Time
	Limit  int
	Offset int
}

type UsageAggregate struct {
	TotalTokens int64
	TotalCost   float64
}

type UsageAggregateFilter struct {
	UserID    int64
	ProjectID *int64
	Start     *time.Time
	End       *time.Time
}

func (r *UsageRepo) ListAdmin(ctx context.Context, filter AdminUsageListFilter) ([]model.UsageRecord, error) {
	query := r.db.WithContext(ctx).Model(&model.UsageRecord{})
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.Start != nil {
		query = query.Where("created_at >= ?", *filter.Start)
	}
	if filter.End != nil {
		query = query.Where("created_at < ?", *filter.End)
	}

	var records []model.UsageRecord
	if err := query.Order("created_at DESC").Limit(filter.Limit).Offset(filter.Offset).Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func (r *UsageRepo) CountAdmin(ctx context.Context, filter AdminUsageListFilter) (int64, error) {
	query := r.db.WithContext(ctx).Model(&model.UsageRecord{})
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.Start != nil {
		query = query.Where("created_at >= ?", *filter.Start)
	}
	if filter.End != nil {
		query = query.Where("created_at < ?", *filter.End)
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *UsageRepo) CountByScope(ctx context.Context, filter UsageAggregateFilter) (int64, error) {
	query := r.db.WithContext(ctx).
		Model(&model.UsageRecord{}).
		Where("user_id = ?", filter.UserID)
	if filter.ProjectID != nil {
		query = query.Where("project_id = ?", *filter.ProjectID)
	}
	if filter.Start != nil {
		query = query.Where("created_at >= ?", *filter.Start)
	}
	if filter.End != nil {
		query = query.Where("created_at < ?", *filter.End)
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *UsageRepo) AggregateByScope(ctx context.Context, filter UsageAggregateFilter) (UsageAggregate, error) {
	query := r.db.WithContext(ctx).
		Model(&model.UsageRecord{}).
		Where("user_id = ?", filter.UserID)
	if filter.ProjectID != nil {
		query = query.Where("project_id = ?", *filter.ProjectID)
	}
	if filter.Start != nil {
		query = query.Where("created_at >= ?", *filter.Start)
	}
	if filter.End != nil {
		query = query.Where("created_at < ?", *filter.End)
	}
	var result UsageAggregate
	if err := query.Select("COALESCE(SUM(total_tokens), 0) AS total_tokens, COALESCE(SUM(cost), 0) AS total_cost").Scan(&result).Error; err != nil {
		return UsageAggregate{}, err
	}
	return result, nil
}
