package repo

import (
	"context"
	"errors"
	"time"

	"deepspace/internal/model"

	"gorm.io/gorm"
)

type PlanUsageRepo struct {
	db *gorm.DB
}

func NewPlanUsageRepo(db *gorm.DB) *PlanUsageRepo {
	return &PlanUsageRepo{db: db}
}

func (r *PlanUsageRepo) GetBySubscriptionPeriod(ctx context.Context, subscriptionID int64, periodStart time.Time, periodEnd *time.Time) (*model.PlanUsage, error) {
	var item model.PlanUsage
	query := r.db.WithContext(ctx).
		Where("subscription_id = ?", subscriptionID).
		Where("period_start = ?", periodStart)
	if periodEnd == nil {
		query = query.Where("period_end IS NULL")
	} else {
		query = query.Where("period_end = ?", *periodEnd)
	}
	err := query.First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (r *PlanUsageRepo) AddUsage(ctx context.Context, subscriptionID, userID int64, periodStart time.Time, periodEnd *time.Time, tokenDelta, requestDelta int64) (*model.PlanUsage, error) {
	if tokenDelta == 0 && requestDelta == 0 {
		return r.GetBySubscriptionPeriod(ctx, subscriptionID, periodStart, periodEnd)
	}

	item, err := r.GetBySubscriptionPeriod(ctx, subscriptionID, periodStart, periodEnd)
	if err != nil {
		return nil, err
	}
	if item == nil {
		item = &model.PlanUsage{
			SubscriptionID: subscriptionID,
			UserID:         userID,
			PeriodStart:    periodStart,
			PeriodEnd:      periodEnd,
			UsedTokens:     tokenDelta,
			UsedRequests:   requestDelta,
		}
		if err := r.db.WithContext(ctx).Create(item).Error; err != nil {
			return nil, err
		}
		return item, nil
	}

	updates := map[string]any{}
	if tokenDelta != 0 {
		updates["used_tokens"] = gorm.Expr("used_tokens + ?", tokenDelta)
	}
	if requestDelta != 0 {
		updates["used_requests"] = gorm.Expr("used_requests + ?", requestDelta)
	}
	if err := r.db.WithContext(ctx).Model(&model.PlanUsage{}).Where("id = ?", item.ID).Updates(updates).Error; err != nil {
		return nil, err
	}

	return r.GetBySubscriptionPeriod(ctx, subscriptionID, periodStart, periodEnd)
}
