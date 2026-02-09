package plan

import (
	"context"
	"errors"
	"math"
	"strings"
	"time"

	"deepspace/internal/model"
	"deepspace/internal/repo"
)

var (
	ErrInvalidPlanName          = errors.New("invalid plan name")
	ErrInvalidPlanQuota         = errors.New("invalid plan quota")
	ErrInvalidPlanCycle         = errors.New("invalid plan cycle")
	ErrInvalidPlanPrice         = errors.New("invalid plan price")
	ErrInvalidPlanCurrency      = errors.New("invalid plan currency")
	ErrPlanNotFound             = errors.New("plan not found")
	ErrInvalidSubscriptionTime  = errors.New("invalid subscription time")
	ErrActiveSubscriptionExists = errors.New("active subscription exists")
)

type Service struct {
	planRepo         *repo.PlanRepo
	subscriptionRepo *repo.PlanSubscriptionRepo
	usageRepo        *repo.PlanUsageRepo
}

func New(planRepo *repo.PlanRepo, subscriptionRepo *repo.PlanSubscriptionRepo, usageRepo *repo.PlanUsageRepo) *Service {
	return &Service{
		planRepo:         planRepo,
		subscriptionRepo: subscriptionRepo,
		usageRepo:        usageRepo,
	}
}

type PlanCreateInput struct {
	Name              string
	Status            string
	IncludedTokens    int64
	IncludedRequests  int64
	ResetIntervalDays int
	Price             float64
	Currency          string
}

type PlanUpdateInput struct {
	Name              *string
	Status            *string
	IncludedTokens    *int64
	IncludedRequests  *int64
	ResetIntervalDays *int
	Price             *float64
	Currency          *string
}

type SubscriptionCreateInput struct {
	UserID  int64
	PlanID  int64
	Status  string
	StartAt time.Time
	EndAt   *time.Time
}

type SubscriptionUpdateInput struct {
	Status  *string
	StartAt *time.Time
	EndAt   *time.Time
}

type ActivePlanQuota struct {
	PlanID         int64
	SubscriptionID int64
	Type           string
	Included       int64
	Used           int64
	Remaining      int64
	PeriodStart    time.Time
	PeriodEnd      *time.Time
}

func (s *Service) CreatePlan(ctx context.Context, input PlanCreateInput) (*model.Plan, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, ErrInvalidPlanName
	}
	if !isValidQuota(input.IncludedTokens, input.IncludedRequests) {
		return nil, ErrInvalidPlanQuota
	}
	resetInterval := normalizeResetInterval(input.ResetIntervalDays)
	if resetInterval == 0 {
		return nil, ErrInvalidPlanCycle
	}
	if input.Price < 0 {
		return nil, ErrInvalidPlanPrice
	}
	currency := normalizeCurrency(input.Currency)
	if currency == "" {
		return nil, ErrInvalidPlanCurrency
	}
	status := normalizeStatus(input.Status)
	if status == "" {
		status = "active"
	}

	plan := &model.Plan{
		Name:              name,
		Status:            status,
		IncludedTokens:    input.IncludedTokens,
		IncludedRequests:  input.IncludedRequests,
		ResetIntervalDays: resetInterval,
		Price:             input.Price,
		Currency:          currency,
	}
	if err := s.planRepo.Create(ctx, plan); err != nil {
		return nil, err
	}
	return plan, nil
}

func (s *Service) UpdatePlan(ctx context.Context, id int64, input PlanUpdateInput) (*model.Plan, error) {
	if id <= 0 {
		return nil, ErrPlanNotFound
	}
	plan, err := s.planRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if plan == nil {
		return nil, ErrPlanNotFound
	}
	updates := map[string]any{}
	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return nil, ErrInvalidPlanName
		}
		updates["name"] = name
	}
	if input.Status != nil {
		status := normalizeStatus(*input.Status)
		if status == "" {
			return nil, ErrInvalidPlanName
		}
		updates["status"] = status
	}
	if input.IncludedTokens != nil {
		updates["included_tokens"] = *input.IncludedTokens
	}
	if input.IncludedRequests != nil {
		updates["included_requests"] = *input.IncludedRequests
	}
	if input.IncludedTokens != nil || input.IncludedRequests != nil {
		finalTokens := plan.IncludedTokens
		finalRequests := plan.IncludedRequests
		if input.IncludedTokens != nil {
			finalTokens = *input.IncludedTokens
		}
		if input.IncludedRequests != nil {
			finalRequests = *input.IncludedRequests
		}
		if !isValidQuota(finalTokens, finalRequests) {
			return nil, ErrInvalidPlanQuota
		}
	}
	if input.ResetIntervalDays != nil {
		if normalizeResetInterval(*input.ResetIntervalDays) == 0 {
			return nil, ErrInvalidPlanCycle
		}
		updates["reset_interval_days"] = *input.ResetIntervalDays
	}
	if input.Price != nil {
		if *input.Price < 0 {
			return nil, ErrInvalidPlanPrice
		}
		updates["price"] = *input.Price
	}
	if input.Currency != nil {
		currency := normalizeCurrency(*input.Currency)
		if currency == "" {
			return nil, ErrInvalidPlanCurrency
		}
		updates["currency"] = currency
	}
	return s.planRepo.Update(ctx, id, updates)
}

func (s *Service) ListPlans(ctx context.Context) ([]model.Plan, error) {
	return s.planRepo.List(ctx)
}

func (s *Service) ListActivePlans(ctx context.Context) ([]model.Plan, error) {
	items, err := s.planRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]model.Plan, 0, len(items))
	for _, item := range items {
		if strings.ToLower(strings.TrimSpace(item.Status)) == "active" {
			result = append(result, item)
		}
	}
	return result, nil
}

func (s *Service) CreateSubscription(ctx context.Context, input SubscriptionCreateInput) (*model.PlanSubscription, error) {
	if input.UserID <= 0 || input.PlanID <= 0 {
		return nil, ErrPlanNotFound
	}
	if input.EndAt != nil && input.EndAt.Before(input.StartAt) {
		return nil, ErrInvalidSubscriptionTime
	}
	status := normalizeStatus(input.Status)
	if status == "" {
		status = "active"
	}
	if status == "active" {
		if existing, err := s.subscriptionRepo.GetActiveByOrg(ctx, input.UserID, time.Now().UTC()); err != nil {
			return nil, err
		} else if existing != nil {
			return nil, ErrActiveSubscriptionExists
		}
	}
	item := &model.PlanSubscription{
		UserID:  input.UserID,
		PlanID:  input.PlanID,
		Status:  status,
		StartAt: input.StartAt,
		EndAt:   input.EndAt,
	}
	if err := s.subscriptionRepo.Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) UpdateSubscription(ctx context.Context, id int64, input SubscriptionUpdateInput) (*model.PlanSubscription, error) {
	if id <= 0 {
		return nil, ErrPlanNotFound
	}
	updates := map[string]any{}
	if input.Status != nil {
		status := normalizeStatus(*input.Status)
		if status == "" {
			return nil, ErrInvalidPlanName
		}
		updates["status"] = status
	}
	if input.StartAt != nil {
		updates["start_at"] = *input.StartAt
	}
	if input.EndAt != nil {
		updates["end_at"] = *input.EndAt
	}
	return s.subscriptionRepo.Update(ctx, id, updates)
}

func (s *Service) GetActiveSubscription(ctx context.Context, userID int64, now time.Time) (*model.PlanSubscription, error) {
	if userID <= 0 {
		return nil, nil
	}
	return s.subscriptionRepo.GetActiveByOrg(ctx, userID, now)
}

func (s *Service) GetActivePlanQuota(ctx context.Context, userID int64, now time.Time) (*ActivePlanQuota, bool, error) {
	subscription, err := s.subscriptionRepo.GetActiveByOrg(ctx, userID, now)
	if err != nil {
		return nil, false, err
	}
	if subscription == nil {
		return nil, false, nil
	}
	plan, err := s.planRepo.GetByID(ctx, subscription.PlanID)
	if err != nil {
		return nil, false, err
	}
	if plan == nil {
		return nil, false, ErrPlanNotFound
	}
	quotaType, included := resolveQuota(plan)
	if quotaType == "" || included <= 0 {
		return nil, false, nil
	}
	periodStart, periodEnd := resolveUsagePeriod(subscription.StartAt, subscription.EndAt, plan.ResetIntervalDays, now)
	usageItem, err := s.usageRepo.GetBySubscriptionPeriod(ctx, subscription.ID, periodStart, periodEnd)
	if err != nil {
		return nil, false, err
	}
	var used int64
	if usageItem != nil {
		if quotaType == "token" {
			used = usageItem.UsedTokens
		} else {
			used = usageItem.UsedRequests
		}
	}
	remaining := included - used
	if remaining < 0 {
		remaining = 0
	}
	return &ActivePlanQuota{
		PlanID:         plan.ID,
		SubscriptionID: subscription.ID,
		Type:           quotaType,
		Included:       included,
		Used:           used,
		Remaining:      remaining,
		PeriodStart:    periodStart,
		PeriodEnd:      periodEnd,
	}, true, nil
}

type QuotaApplyResult struct {
	Applied     bool
	Type        string
	Included    int64
	UsedBefore  int64
	UsedAfter   int64
	Overage     int64
	Remaining   int64
	PeriodStart time.Time
	PeriodEnd   *time.Time
}

func (s *Service) ApplyQuota(ctx context.Context, userID int64, now time.Time, tokenUnits int64, requestUnits int64) (*QuotaApplyResult, error) {
	quota, ok, err := s.GetActivePlanQuota(ctx, userID, now)
	if err != nil || !ok || quota == nil {
		return &QuotaApplyResult{Applied: false}, err
	}
	var units int64
	if quota.Type == "token" {
		units = tokenUnits
	} else {
		units = requestUnits
	}
	if units < 0 {
		units = 0
	}
	usedBefore := quota.Used
	remainingBefore := quota.Remaining
	overage := int64(0)
	if units > remainingBefore {
		overage = units - remainingBefore
	}
	usedAfter := usedBefore + units
	remainingAfter := quota.Included - usedAfter
	if remainingAfter < 0 {
		remainingAfter = 0
	}
	if units > 0 {
		_, err = s.usageRepo.AddUsage(ctx, quota.SubscriptionID, userID, quota.PeriodStart, quota.PeriodEnd, applyTokenDelta(quota.Type, units), applyRequestDelta(quota.Type, units))
		if err != nil {
			return nil, err
		}
	}
	return &QuotaApplyResult{
		Applied:     true,
		Type:        quota.Type,
		Included:    quota.Included,
		UsedBefore:  usedBefore,
		UsedAfter:   usedAfter,
		Overage:     overage,
		Remaining:   remainingAfter,
		PeriodStart: quota.PeriodStart,
		PeriodEnd:   quota.PeriodEnd,
	}, nil
}

func normalizeStatus(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" {
		return ""
	}
	switch value {
	case "active", "disabled", "expired", "canceled":
		return value
	default:
		return ""
	}
}

func normalizeResetInterval(value int) int {
	if value == 0 {
		return 30
	}
	if value < 1 || value > 365 {
		return 0
	}
	return value
}

func normalizeCurrency(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "CNY"
	}
	return strings.ToUpper(value)
}

func isValidQuota(tokens, requests int64) bool {
	if tokens > 0 && requests > 0 {
		return false
	}
	if tokens <= 0 && requests <= 0 {
		return false
	}
	return true
}

func resolveQuota(plan *model.Plan) (string, int64) {
	if plan == nil {
		return "", 0
	}
	if plan.IncludedTokens > 0 && plan.IncludedRequests == 0 {
		return "token", plan.IncludedTokens
	}
	if plan.IncludedRequests > 0 && plan.IncludedTokens == 0 {
		return "request", plan.IncludedRequests
	}
	return "", 0
}

func resolveUsagePeriod(start time.Time, end *time.Time, intervalDays int, now time.Time) (time.Time, *time.Time) {
	intervalDays = normalizeResetInterval(intervalDays)
	if intervalDays == 0 {
		intervalDays = 30
	}
	startUTC := start.UTC()
	nowUTC := now.UTC()
	if nowUTC.Before(startUTC) {
		periodStart := startUTC
		periodEnd := startUTC.AddDate(0, 0, intervalDays)
		periodEnd = clampPeriodEnd(periodEnd, end)
		return periodStart, &periodEnd
	}

	elapsedDays := int(math.Floor(nowUTC.Sub(startUTC).Hours() / 24))
	periodIndex := elapsedDays / intervalDays
	periodStart := startUTC.AddDate(0, 0, periodIndex*intervalDays)
	periodEnd := periodStart.AddDate(0, 0, intervalDays)
	periodEnd = clampPeriodEnd(periodEnd, end)
	return periodStart, &periodEnd
}

func clampPeriodEnd(periodEnd time.Time, end *time.Time) time.Time {
	if end == nil {
		return periodEnd
	}
	endUTC := end.UTC()
	if periodEnd.After(endUTC) {
		return endUTC
	}
	return periodEnd
}

func applyTokenDelta(quotaType string, units int64) int64 {
	if quotaType == "token" {
		return units
	}
	return 0
}

func applyRequestDelta(quotaType string, units int64) int64 {
	if quotaType == "request" {
		return units
	}
	return 0
}
