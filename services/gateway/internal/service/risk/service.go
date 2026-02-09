package risk

import (
	"context"
	"errors"
	"strings"

	"deepspace/internal/model"
	"deepspace/internal/repo"
)

var (
	ErrInvalidScope  = errors.New("invalid scope")
	ErrInvalidStatus = errors.New("invalid status")
	ErrInvalidName   = errors.New("invalid name")
	ErrInvalidRule   = errors.New("invalid rule")
	ErrInvalidCycle  = errors.New("invalid cycle")
	ErrInvalidPolicy = errors.New("invalid policy")
	ErrInvalidIPRule = errors.New("invalid ip rule")
)

const (
	ScopeGlobal  = "global"
	ScopeUser    = "user"
	ScopeProject = "project"
)

type Service struct {
	policyRepo *repo.RiskPolicyRepo
	rateRepo   *repo.RateLimitRepo
	ipRepo     *repo.IPRuleRepo
	budgetRepo *repo.BudgetCapRepo
}

func New(policyRepo *repo.RiskPolicyRepo, rateRepo *repo.RateLimitRepo, ipRepo *repo.IPRuleRepo, budgetRepo *repo.BudgetCapRepo) *Service {
	return &Service{policyRepo: policyRepo, rateRepo: rateRepo, ipRepo: ipRepo, budgetRepo: budgetRepo}
}

type PolicyCreateInput struct {
	Name      string
	Scope     string
	UserID    *int64
	ProjectID *int64
	Status    string
	Priority  int
}

type PolicyUpdateInput struct {
	Name      *string
	Scope     *string
	UserID    *int64
	ProjectID *int64
	Status    *string
	Priority  *int
}

type RateLimitInput struct {
	PolicyID      int64
	WindowSeconds int
	MaxRequests   int
	MaxTokens     int
	Status        string
}

type RateLimitUpdateInput struct {
	WindowSeconds *int
	MaxRequests   *int
	MaxTokens     *int
	Status        *string
}

type IPRuleInput struct {
	PolicyID int64
	Type     string
	IP       *string
	CIDR     *string
	Status   string
}

type IPRuleUpdateInput struct {
	Type   *string
	IP     *string
	CIDR   *string
	Status *string
}

type BudgetCapInput struct {
	PolicyID int64
	Cycle    string
	MaxCost  float64
	Currency string
	Status   string
}

type BudgetCapUpdateInput struct {
	Cycle    *string
	MaxCost  *float64
	Currency *string
	Status   *string
}

func (s *Service) CreatePolicy(ctx context.Context, input PolicyCreateInput) (*model.RiskPolicy, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, ErrInvalidName
	}
	scope := normalizeScope(input.Scope)
	if scope == "" {
		return nil, ErrInvalidScope
	}
	if !validScopeTarget(scope, input.UserID, input.ProjectID) {
		return nil, ErrInvalidScope
	}
	status := normalizeStatus(input.Status)
	if status == "" {
		status = "active"
	}
	item := &model.RiskPolicy{
		Name:      name,
		Scope:     scope,
		UserID:    input.UserID,
		ProjectID: input.ProjectID,
		Status:    status,
		Priority:  input.Priority,
	}
	if err := s.policyRepo.Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) UpdatePolicy(ctx context.Context, id int64, input PolicyUpdateInput) (*model.RiskPolicy, error) {
	if id <= 0 {
		return nil, ErrInvalidPolicy
	}
	item, err := s.policyRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrInvalidPolicy
	}
	updates := map[string]any{}
	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return nil, ErrInvalidName
		}
		updates["name"] = name
	}
	var scope string
	var userID *int64
	var projectID *int64
	if input.Scope != nil {
		scope = normalizeScope(*input.Scope)
		if scope == "" {
			return nil, ErrInvalidScope
		}
		updates["scope"] = scope
	} else {
		scope = item.Scope
	}
	if input.UserID != nil {
		userID = input.UserID
		updates["user_id"] = input.UserID
	} else {
		userID = item.UserID
	}
	if input.ProjectID != nil {
		projectID = input.ProjectID
		updates["project_id"] = input.ProjectID
	} else {
		projectID = item.ProjectID
	}
	if !validScopeTarget(scope, userID, projectID) {
		return nil, ErrInvalidScope
	}
	if input.Status != nil {
		status := normalizeStatus(*input.Status)
		if status == "" {
			return nil, ErrInvalidStatus
		}
		updates["status"] = status
	}
	if input.Priority != nil {
		updates["priority"] = *input.Priority
	}
	return s.policyRepo.Update(ctx, id, updates)
}

func (s *Service) DeletePolicy(ctx context.Context, id int64) (bool, error) {
	if id <= 0 {
		return false, ErrInvalidPolicy
	}
	return s.policyRepo.Delete(ctx, id)
}

func (s *Service) ListPolicies(ctx context.Context, filter repo.RiskPolicyFilter) ([]model.RiskPolicy, int64, error) {
	return s.policyRepo.List(ctx, filter)
}

func (s *Service) CreateRateLimit(ctx context.Context, input RateLimitInput) (*model.RateLimit, error) {
	if input.PolicyID <= 0 {
		return nil, ErrInvalidPolicy
	}
	if input.WindowSeconds <= 0 || (input.MaxRequests <= 0 && input.MaxTokens <= 0) {
		return nil, ErrInvalidRule
	}
	status := normalizeStatus(input.Status)
	if status == "" {
		status = "active"
	}
	item := &model.RateLimit{
		PolicyID:      input.PolicyID,
		WindowSeconds: input.WindowSeconds,
		MaxRequests:   input.MaxRequests,
		MaxTokens:     input.MaxTokens,
		Status:        status,
	}
	if err := s.rateRepo.Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) UpdateRateLimit(ctx context.Context, id int64, input RateLimitUpdateInput) (*model.RateLimit, error) {
	if id <= 0 {
		return nil, ErrInvalidRule
	}
	updates := map[string]any{}
	if input.WindowSeconds != nil {
		if *input.WindowSeconds <= 0 {
			return nil, ErrInvalidRule
		}
		updates["window_seconds"] = *input.WindowSeconds
	}
	if input.MaxRequests != nil {
		updates["max_requests"] = *input.MaxRequests
	}
	if input.MaxTokens != nil {
		updates["max_tokens"] = *input.MaxTokens
	}
	if input.Status != nil {
		status := normalizeStatus(*input.Status)
		if status == "" {
			return nil, ErrInvalidStatus
		}
		updates["status"] = status
	}
	return s.rateRepo.Update(ctx, id, updates)
}

func (s *Service) DeleteRateLimit(ctx context.Context, id int64) (bool, error) {
	if id <= 0 {
		return false, ErrInvalidRule
	}
	return s.rateRepo.Delete(ctx, id)
}

func (s *Service) ListRateLimits(ctx context.Context, filter repo.RateLimitFilter) ([]model.RateLimit, int64, error) {
	return s.rateRepo.List(ctx, filter)
}

func (s *Service) CreateIPRule(ctx context.Context, input IPRuleInput) (*model.IPRule, error) {
	if input.PolicyID <= 0 {
		return nil, ErrInvalidPolicy
	}
	typeValue := strings.ToLower(strings.TrimSpace(input.Type))
	if typeValue != "allow" && typeValue != "deny" {
		return nil, ErrInvalidIPRule
	}
	if (input.IP == nil || strings.TrimSpace(*input.IP) == "") && (input.CIDR == nil || strings.TrimSpace(*input.CIDR) == "") {
		return nil, ErrInvalidIPRule
	}
	status := normalizeStatus(input.Status)
	if status == "" {
		status = "active"
	}
	item := &model.IPRule{
		PolicyID: input.PolicyID,
		Type:     typeValue,
		IP:       input.IP,
		CIDR:     input.CIDR,
		Status:   status,
	}
	if err := s.ipRepo.Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) UpdateIPRule(ctx context.Context, id int64, input IPRuleUpdateInput) (*model.IPRule, error) {
	if id <= 0 {
		return nil, ErrInvalidIPRule
	}
	updates := map[string]any{}
	if input.Type != nil {
		value := strings.ToLower(strings.TrimSpace(*input.Type))
		if value != "allow" && value != "deny" {
			return nil, ErrInvalidIPRule
		}
		updates["type"] = value
	}
	if input.IP != nil {
		if strings.TrimSpace(*input.IP) == "" {
			return nil, ErrInvalidIPRule
		}
		updates["ip"] = input.IP
	}
	if input.CIDR != nil {
		if strings.TrimSpace(*input.CIDR) == "" {
			return nil, ErrInvalidIPRule
		}
		updates["cidr"] = input.CIDR
	}
	if input.Status != nil {
		status := normalizeStatus(*input.Status)
		if status == "" {
			return nil, ErrInvalidStatus
		}
		updates["status"] = status
	}
	return s.ipRepo.Update(ctx, id, updates)
}

func (s *Service) DeleteIPRule(ctx context.Context, id int64) (bool, error) {
	if id <= 0 {
		return false, ErrInvalidIPRule
	}
	return s.ipRepo.Delete(ctx, id)
}

func (s *Service) ListIPRules(ctx context.Context, filter repo.IPRuleFilter) ([]model.IPRule, int64, error) {
	return s.ipRepo.List(ctx, filter)
}

func (s *Service) CreateBudgetCap(ctx context.Context, input BudgetCapInput) (*model.BudgetCap, error) {
	if input.PolicyID <= 0 {
		return nil, ErrInvalidPolicy
	}
	cycle := normalizeCycle(input.Cycle)
	if cycle == "" {
		return nil, ErrInvalidCycle
	}
	if input.MaxCost <= 0 {
		return nil, ErrInvalidRule
	}
	currency := normalizeCurrency(input.Currency)
	if currency == "" {
		return nil, ErrInvalidRule
	}
	status := normalizeStatus(input.Status)
	if status == "" {
		status = "active"
	}
	item := &model.BudgetCap{
		PolicyID: input.PolicyID,
		Cycle:    cycle,
		MaxCost:  input.MaxCost,
		Currency: currency,
		Status:   status,
	}
	if err := s.budgetRepo.Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) UpdateBudgetCap(ctx context.Context, id int64, input BudgetCapUpdateInput) (*model.BudgetCap, error) {
	if id <= 0 {
		return nil, ErrInvalidRule
	}
	updates := map[string]any{}
	if input.Cycle != nil {
		cycle := normalizeCycle(*input.Cycle)
		if cycle == "" {
			return nil, ErrInvalidCycle
		}
		updates["cycle"] = cycle
	}
	if input.MaxCost != nil {
		if *input.MaxCost <= 0 {
			return nil, ErrInvalidRule
		}
		updates["max_cost"] = *input.MaxCost
	}
	if input.Currency != nil {
		currency := normalizeCurrency(*input.Currency)
		if currency == "" {
			return nil, ErrInvalidRule
		}
		updates["currency"] = currency
	}
	if input.Status != nil {
		status := normalizeStatus(*input.Status)
		if status == "" {
			return nil, ErrInvalidStatus
		}
		updates["status"] = status
	}
	return s.budgetRepo.Update(ctx, id, updates)
}

func (s *Service) DeleteBudgetCap(ctx context.Context, id int64) (bool, error) {
	if id <= 0 {
		return false, ErrInvalidRule
	}
	return s.budgetRepo.Delete(ctx, id)
}

func (s *Service) ListBudgetCaps(ctx context.Context, filter repo.BudgetCapFilter) ([]model.BudgetCap, int64, error) {
	return s.budgetRepo.List(ctx, filter)
}

func normalizeScope(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case ScopeGlobal, ScopeUser, ScopeProject:
		return strings.ToLower(strings.TrimSpace(value))
	default:
		return ""
	}
}

func normalizeStatus(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "active", "disabled":
		return strings.ToLower(strings.TrimSpace(value))
	default:
		return ""
	}
}

func normalizeCycle(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "daily", "weekly", "monthly":
		return strings.ToLower(strings.TrimSpace(value))
	default:
		return ""
	}
}

func normalizeCurrency(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "CNY"
	}
	return strings.ToUpper(value)
}

func validScopeTarget(scope string, userID, projectID *int64) bool {
	switch scope {
	case ScopeGlobal:
		return userID == nil && projectID == nil
	case ScopeUser:
		return userID != nil && projectID == nil
	case ScopeProject:
		return projectID != nil
	default:
		return false
	}
}
