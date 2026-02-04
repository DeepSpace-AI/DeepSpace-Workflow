package pipeline

import (
	"context"
	"errors"
)

var ErrHalt = errors.New("pipeline halted")

type Step interface {
	Name() string
	Run(ctx context.Context, state *State) error
}

type Pipeline struct {
	steps []Step
}

func New(steps ...Step) *Pipeline {
	return &Pipeline{steps: steps}
}

func (p *Pipeline) Run(ctx context.Context, state *State) error {
	for _, step := range p.steps {
		if state.Halted {
			return ErrHalt
		}
		if err := step.Run(ctx, state); err != nil {
			return err
		}
	}
	return nil
}
