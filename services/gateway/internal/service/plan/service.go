package plan

import (
	"context"
	"errors"
	"strings"
	"time"

	"deepspace/internal/model"
	"deepspace/internal/repo"
)

const (
	BillingModeToken   = "token"
	BillingModeRequest = "request"
)

var (
	ErrInvalidPlanName          = errors.New("invalid plan name")
	ErrInvalidBillingMode       = errors.New("invalid billing mode")
	ErrInvalidPlanPrice         = errors.New("invalid plan price")
	ErrInvalidPlanCurrency      = errors.New("invalid plan currency")
	ErrPlanNotFound             = errors.New("plan not found")
	ErrInvalidSubscriptionTime  = errors.New("invalid subscription time")
	ErrActiveSubscriptionExists = errors.New("active subscription exists")
)

type Service struct {
	planRepo         *repo.PlanRepo
	priceRepo        *repo.PlanModelPriceRepo
	subscriptionRepo *repo.PlanSubscriptionRepo
}

func New(planRepo *repo.PlanRepo, priceRepo *repo.PlanModelPriceRepo, subscriptionRepo *repo.PlanSubscriptionRepo) *Service {
	return &Service{
		planRepo:         planRepo,
		priceRepo:        priceRepo,
		subscriptionRepo: subscriptionRepo,
	}
}

type PlanCreateInput struct {
	Name         string
	Status       string
	Currency     string
	BillingMode  string
	PriceInput   float64
	PriceOutput  float64
	PriceRequest float64
}

type PlanUpdateInput struct {
	Name         *string
	Status       *string
	Currency     *string
	BillingMode  *string
	PriceInput   *float64
	PriceOutput  *float64
	PriceRequest *float64
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

type ActivePlanPrice struct {
	PlanID       int64
	BillingMode  string
	Currency     string
	PriceInput   float64
	PriceOutput  float64
	PriceRequest float64
}

func (s *Service) CreatePlan(ctx context.Context, input PlanCreateInput) (*model.Plan, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, ErrInvalidPlanName
	}
	billingMode := normalizeBillingMode(input.BillingMode)
	if billingMode == "" {
		return nil, ErrInvalidBillingMode
	}
	currency := normalizeCurrency(input.Currency)
	if currency == "" {
		return nil, ErrInvalidPlanCurrency
	}
	if input.PriceInput < 0 || input.PriceOutput < 0 || input.PriceRequest < 0 {
		return nil, ErrInvalidPlanPrice
	}
	status := normalizeStatus(input.Status)
	if status == "" {
		status = "active"
	}

	plan := &model.Plan{
		Name:         name,
		Status:       status,
		Currency:     currency,
		BillingMode:  billingMode,
		PriceInput:   input.PriceInput,
		PriceOutput:  input.PriceOutput,
		PriceRequest: input.PriceRequest,
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
	if input.Currency != nil {
		currency := normalizeCurrency(*input.Currency)
		if currency == "" {
			return nil, ErrInvalidPlanCurrency
		}
		updates["currency"] = currency
	}
	if input.BillingMode != nil {
		billingMode := normalizeBillingMode(*input.BillingMode)
		if billingMode == "" {
			return nil, ErrInvalidBillingMode
		}
		updates["billing_mode"] = billingMode
	}
	if input.PriceInput != nil {
		if *input.PriceInput < 0 {
			return nil, ErrInvalidPlanPrice
		}
		updates["price_input"] = *input.PriceInput
	}
	if input.PriceOutput != nil {
		if *input.PriceOutput < 0 {
			return nil, ErrInvalidPlanPrice
		}
		updates["price_output"] = *input.PriceOutput
	}
	if input.PriceRequest != nil {
		if *input.PriceRequest < 0 {
			return nil, ErrInvalidPlanPrice
		}
		updates["price_request"] = *input.PriceRequest
	}
	return s.planRepo.Update(ctx, id, updates)
}

func (s *Service) ListPlans(ctx context.Context) ([]model.Plan, error) {
	return s.planRepo.List(ctx)
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

func (s *Service) GetActivePlanPrice(ctx context.Context, userID int64, modelID string, now time.Time) (*ActivePlanPrice, bool, error) {
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
	result := &ActivePlanPrice{
		PlanID:       plan.ID,
		BillingMode:  plan.BillingMode,
		Currency:     plan.Currency,
		PriceInput:   plan.PriceInput,
		PriceOutput:  plan.PriceOutput,
		PriceRequest: plan.PriceRequest,
	}
	if modelID == "" {
		return result, true, nil
	}
	override, err := s.priceRepo.GetByPlanModel(ctx, plan.ID, modelID)
	if err != nil {
		return nil, false, err
	}
	if override == nil {
		return result, true, nil
	}
	result.Currency = override.Currency
	result.PriceInput = override.PriceInput
	result.PriceOutput = override.PriceOutput
	result.PriceRequest = override.PriceRequest
	return result, true, nil
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

func normalizeCurrency(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "USD"
	}
	return strings.ToUpper(value)
}

func normalizeBillingMode(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	switch value {
	case BillingModeToken, BillingModeRequest:
		return value
	default:
		return ""
	}
}
