package saga

import (
	"context"

	"github.com/pkg/errors"
)

type Coordinator interface {
	Execute(ctx context.Context, steps ...*Step) error
}

func NewCoordinator() Coordinator {
	return &coordinator{}
}

type coordinator struct{}

func (s coordinator) compensate(ctx context.Context, steps []*Step) error {
	for i := len(steps) - 1; i >= 0; i-- {
		if err := steps[i].Compensate(ctx); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (s coordinator) Execute(ctx context.Context, steps ...*Step) error {
	for i := 0; i < len(steps); i++ {
		if err := steps[i].Handle(ctx); err != nil {
			if err := s.compensate(ctx, steps[:i]); err != nil {
				return errors.WithStack(err)
			}

			return errors.WithStack(err)
		}
	}

	return nil
}
