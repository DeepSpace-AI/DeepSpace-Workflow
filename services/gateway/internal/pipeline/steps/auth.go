package steps

import (
	"context"

	"deepspace/internal/pipeline"
)

type Auth struct{}

func NewAuth() *Auth {
	return &Auth{}
}

func (s *Auth) Name() string {
	return "auth"
}

func (s *Auth) Run(ctx context.Context, state *pipeline.State) error {
	// Auth is currently handled by middleware; this step is a placeholder.
	return nil
}
