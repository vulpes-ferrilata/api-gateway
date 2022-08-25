package interceptors

import (
	"context"
	"strings"

	httpext "github.com/go-playground/pkg/v5/net/http"
	"github.com/pkg/errors"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/context_values"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func NewLocaleInterceptor() *LocaleInterceptor {
	return &LocaleInterceptor{}
}

type LocaleInterceptor struct{}

func (l LocaleInterceptor) ClientUnaryInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	locales := context_values.GetLocales(ctx)

	md := metadata.MD{
		strings.ToLower(httpext.AcceptedLanguage): locales,
	}

	ctx = metadata.NewOutgoingContext(ctx, md)

	if err := invoker(ctx, method, req, reply, cc, opts...); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
