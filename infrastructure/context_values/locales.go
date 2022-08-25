package context_values

import (
	"context"
)

type localesContextKey struct{}

func WithLocales(ctx context.Context, locales []string) context.Context {
	return context.WithValue(ctx, localesContextKey{}, locales)
}

func GetLocales(ctx context.Context) []string {
	locales, _ := ctx.Value(localesContextKey{}).([]string)
	return locales
}
