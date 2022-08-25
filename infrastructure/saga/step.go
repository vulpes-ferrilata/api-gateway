package saga

import "context"

type HandlerFunc func(ctx context.Context) error

type Step struct {
	Handle     HandlerFunc
	Compensate HandlerFunc
}
