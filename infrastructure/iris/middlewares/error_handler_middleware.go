package middlewares

import (
	"fmt"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/hero"
	"github.com/pkg/errors"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/app_errors"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/context_values"
	"google.golang.org/grpc/status"
)

func NewErrorHandlerMiddleware(universalTranslator *ut.UniversalTranslator) *ErrorHandlerMiddleware {
	return &ErrorHandlerMiddleware{
		universalTranslator: universalTranslator,
	}
}

type ErrorHandlerMiddleware struct {
	universalTranslator *ut.UniversalTranslator
}

func (e ErrorHandlerMiddleware) Serve() hero.ErrorHandlerFunc {
	return func(ctx iris.Context, err error) {
		locales := context_values.GetLocales(ctx.Request().Context())
		translator, _ := e.universalTranslator.FindTranslator(locales...)

		if status, ok := status.FromError(errors.Cause(err)); ok {
			err = app_errors.NewStatusError(status)
		}

		if validationErrors, ok := errors.Cause(err).(validator.ValidationErrors); ok {
			err = app_errors.NewRequestValidationError(validationErrors)
		}

		if webErr, ok := errors.Cause(err).(app_errors.WebError); ok {
			problem := webErr.Problem(translator)
			ctx.Problem(problem)
			ctx.StopExecution()
			return
		}

		problem := iris.NewProblem()
		problem.Status(iris.StatusInternalServerError)

		if detail, err := translator.T("internal-error"); err == nil {
			problem.Detail(detail)
		} else {
			problem.Detail("internal-error")
		}

		problem.Key("stacktrace", strings.Split(fmt.Sprintf("%+v", err), "\n"))

		ctx.Problem(problem)
		ctx.StopExecution()
	}
}
