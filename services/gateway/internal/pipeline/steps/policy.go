package steps

import (
	"context"
	"errors"
	"net"
	"strings"
	"time"

	"deepspace/internal/model"
	"deepspace/internal/pipeline"
	"deepspace/internal/repo"
	"deepspace/internal/service/risk"
	"deepspace/internal/service/usage"
)

var (
	ErrRiskIPDenied       = errors.New("risk ip denied")
	ErrRiskRateLimited    = errors.New("risk rate limited")
	ErrRiskBudgetExceeded = errors.New("risk budget exceeded")
)

type Policy struct {
	risk  *risk.Service
	usage *usage.Service
}

func NewPolicy(riskSvc *risk.Service, usageSvc *usage.Service) *Policy {
	return &Policy{risk: riskSvc, usage: usageSvc}
}

func (s *Policy) Name() string {
	return "policy"
}

func (s *Policy) Run(ctx context.Context, state *pipeline.State) error {
	if s == nil || s.risk == nil || state == nil {
		return nil
	}
	if state.UserID <= 0 {
		return nil
	}

	policy, err := s.resolvePolicy(ctx, state)
	if err != nil || policy == nil {
		return err
	}

	if err := s.applyIPRules(ctx, state, policy.ID); err != nil {
		return err
	}
	if err := s.applyRateLimits(ctx, state, policy.ID); err != nil {
		return err
	}
	if err := s.applyBudgetCaps(ctx, state, policy.ID); err != nil {
		return err
	}

	return nil
}

func (s *Policy) resolvePolicy(ctx context.Context, state *pipeline.State) (*model.RiskPolicy, error) {
	const statusActive = "active"

	if state.ProjectID != nil {
		item, err := s.pickPolicy(ctx, repo.RiskPolicyFilter{
			Scope:     risk.ScopeProject,
			ProjectID: state.ProjectID,
			Status:    statusActive,
			Limit:     1,
			Offset:    0,
		})
		if err != nil || item != nil {
			return item, err
		}
	}

	userID := state.UserID
	item, err := s.pickPolicy(ctx, repo.RiskPolicyFilter{
		Scope:  risk.ScopeUser,
		UserID: &userID,
		Status: statusActive,
		Limit:  1,
		Offset: 0,
	})
	if err != nil || item != nil {
		return item, err
	}

	return s.pickPolicy(ctx, repo.RiskPolicyFilter{
		Scope:  risk.ScopeGlobal,
		Status: statusActive,
		Limit:  1,
		Offset: 0,
	})
}

func (s *Policy) pickPolicy(ctx context.Context, filter repo.RiskPolicyFilter) (*model.RiskPolicy, error) {
	items, _, err := s.risk.ListPolicies(ctx, filter)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, nil
	}
	item := items[0]
	return &item, nil
}

func (s *Policy) applyIPRules(ctx context.Context, state *pipeline.State, policyID int64) error {
	policyIDValue := policyID
	items, _, err := s.risk.ListIPRules(ctx, repo.IPRuleFilter{
		PolicyID: &policyIDValue,
		Status:   "active",
		Limit:    1000,
		Offset:   0,
	})
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return nil
	}

	clientIP := net.ParseIP(strings.TrimSpace(getMetaString(state.Meta, "client_ip")))
	if clientIP == nil {
		for _, rule := range items {
			if strings.EqualFold(rule.Type, "allow") {
				return ErrRiskIPDenied
			}
		}
		return nil
	}

	hasAllow := false
	allowMatched := false
	for _, rule := range items {
		if strings.EqualFold(rule.Type, "allow") {
			hasAllow = true
		}
		if matchIPRule(clientIP, rule) {
			if strings.EqualFold(rule.Type, "deny") {
				return ErrRiskIPDenied
			}
			if strings.EqualFold(rule.Type, "allow") {
				allowMatched = true
			}
		}
	}
	if hasAllow && !allowMatched {
		return ErrRiskIPDenied
	}
	return nil
}

func (s *Policy) applyRateLimits(ctx context.Context, state *pipeline.State, policyID int64) error {
	if s.usage == nil {
		return nil
	}
	policyIDValue := policyID
	items, _, err := s.risk.ListRateLimits(ctx, repo.RateLimitFilter{
		PolicyID: &policyIDValue,
		Status:   "active",
		Limit:    1000,
		Offset:   0,
	})
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return nil
	}

	now := time.Now().UTC()
	for _, rule := range items {
		if rule.WindowSeconds <= 0 {
			continue
		}
		windowStart := now.Add(-time.Duration(rule.WindowSeconds) * time.Second)
		if rule.MaxRequests > 0 {
			count, err := s.usage.CountByScope(ctx, usage.AggregateInput{
				UserID:    state.UserID,
				ProjectID: state.ProjectID,
				Start:     &windowStart,
				End:       &now,
			})
			if err != nil {
				return err
			}
			if count+1 > int64(rule.MaxRequests) {
				return ErrRiskRateLimited
			}
		}
		if rule.MaxTokens > 0 {
			agg, err := s.usage.AggregateByScope(ctx, usage.AggregateInput{
				UserID:    state.UserID,
				ProjectID: state.ProjectID,
				Start:     &windowStart,
				End:       &now,
			})
			if err != nil {
				return err
			}
			if agg.TotalTokens >= int64(rule.MaxTokens) {
				return ErrRiskRateLimited
			}
		}
	}

	return nil
}

func (s *Policy) applyBudgetCaps(ctx context.Context, state *pipeline.State, policyID int64) error {
	if s.usage == nil {
		return nil
	}
	policyIDValue := policyID
	items, _, err := s.risk.ListBudgetCaps(ctx, repo.BudgetCapFilter{
		PolicyID: &policyIDValue,
		Status:   "active",
		Limit:    1000,
		Offset:   0,
	})
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return nil
	}

	now := time.Now().UTC()
	metaCurrency := strings.ToUpper(strings.TrimSpace(getMetaString(state.Meta, "currency")))
	for _, cap := range items {
		cycleStart, ok := resolveCycleStart(cap.Cycle, now)
		if !ok {
			continue
		}
		if !currencyMatch(cap.Currency, metaCurrency) {
			continue
		}
		agg, err := s.usage.AggregateByScope(ctx, usage.AggregateInput{
			UserID:    state.UserID,
			ProjectID: state.ProjectID,
			Start:     &cycleStart,
			End:       &now,
		})
		if err != nil {
			return err
		}
		if cap.MaxCost > 0 && agg.TotalCost >= cap.MaxCost {
			return ErrRiskBudgetExceeded
		}
	}

	return nil
}

func matchIPRule(clientIP net.IP, rule model.IPRule) bool {
	if clientIP == nil {
		return false
	}
	if rule.IP != nil {
		parsed := net.ParseIP(strings.TrimSpace(*rule.IP))
		if parsed != nil && parsed.Equal(clientIP) {
			return true
		}
	}
	if rule.CIDR != nil {
		_, cidr, err := net.ParseCIDR(strings.TrimSpace(*rule.CIDR))
		if err == nil && cidr != nil && cidr.Contains(clientIP) {
			return true
		}
	}
	return false
}

func resolveCycleStart(cycle string, now time.Time) (time.Time, bool) {
	now = now.UTC()
	switch strings.ToLower(strings.TrimSpace(cycle)) {
	case "daily":
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC), true
	case "weekly":
		weekday := int(now.Weekday())
		offset := (weekday + 6) % 7
		start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		return start.AddDate(0, 0, -offset), true
	case "monthly":
		return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC), true
	default:
		return time.Time{}, false
	}
}

func currencyMatch(capCurrency, metaCurrency string) bool {
	value := strings.ToUpper(strings.TrimSpace(capCurrency))
	if value == "" {
		return true
	}
	if metaCurrency == "" {
		return true
	}
	return value == metaCurrency
}

func getMetaString(meta map[string]any, key string) string {
	if meta == nil {
		return ""
	}
	value, ok := meta[key]
	if !ok {
		return ""
	}
	if text, ok := value.(string); ok {
		return text
	}
	return ""
}
