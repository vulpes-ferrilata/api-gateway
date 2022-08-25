package middlewares

import (
	"fmt"

	ut "github.com/go-playground/universal-translator"
	"github.com/kataras/iris/v12"
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

func (e ErrorHandlerMiddleware) Handle(ctx iris.Context, err error) {
	if status, ok := status.FromError(errors.Cause(err)); ok {
		err = app_errors.NewStatusError(status)
	}

	if webErr, ok := errors.Cause(err).(app_errors.WebError); ok {
		locales := context_values.GetLocales(ctx.Request().Context())
		translator, _ := e.universalTranslator.FindTranslator(locales...)
		problem := webErr.Problem(translator)
		ctx.Problem(problem)
		ctx.StopExecution()
		return
	}

	problem := iris.NewProblem()
	problem.Status(iris.StatusInternalServerError)
	problem.Detail("something went wrong")
	problem.Key("stackTraces", fmt.Sprintf("%+v", err))
	ctx.Problem(problem)
	ctx.StopExecution()
}
