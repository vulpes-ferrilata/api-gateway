package middlewares

import (
	"github.com/go-playground/pure"
	"github.com/kataras/iris/v12"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/context_values"
)

func NewLocaleMiddleware() *LocaleMiddleware {
	return &LocaleMiddleware{}
}

type LocaleMiddleware struct{}

func (l LocaleMiddleware) Serve() iris.Handler {
	return func(ctx iris.Context) {
		request := ctx.Request()
		requestCtx := request.Context()
		locales := pure.AcceptedLanguages(request)

		requestCtx = context_values.WithLocales(requestCtx, locales)
		request = request.WithContext(requestCtx)
		ctx.ResetRequest(request)

		ctx.Next()
	}
}
