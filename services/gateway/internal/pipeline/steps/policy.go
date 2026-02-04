package steps

import (
	"context"

	"deepspace/internal/pipeline"
)

type Policy struct{}

func NewPolicy() *Policy {
	return &Policy{}
}

func (s *Policy) Name() string {
	return "policy"
}

func (s *Policy) Run(ctx context.Context, state *pipeline.State) error {
	// Policy enforcement placeholder.
	return nil
}
