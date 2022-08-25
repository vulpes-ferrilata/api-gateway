package context_values

import (
	"context"
)

type userIDContextKey struct{}

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDContextKey{}, userID)
}

func GetUserID(ctx context.Context) string {
	return ctx.Value(userIDContextKey{}).(string)
}
