package steps

import (
	"context"

	"deepspace/internal/pipeline"
	"deepspace/internal/service/billing"
)

type BudgetHold struct {
	billing *billing.Service
}

func NewBudgetHold(billingSvc *billing.Service) *BudgetHold {
	return &BudgetHold{billing: billingSvc}
}

func (s *BudgetHold) Name() string {
	return "budget_hold"
}

func (s *BudgetHold) Run(ctx context.Context, state *pipeline.State) error {
	if s.billing == nil {
		return nil
	}
	if state.CostAmount <= 0 || state.RefID == "" {
		return nil
	}
	_, err := s.billing.Hold(ctx, state.UserID, state.CostAmount, state.RefID, map[string]any{"source": "pipeline"})
	return err
}
