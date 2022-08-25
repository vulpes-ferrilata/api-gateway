package app_errors

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/kataras/iris/v12"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewStatusError(status *status.Status) AppError {
	return &statusError{
		status: status,
	}
}

type statusError struct {
	status *status.Status
}

func (s statusError) Error() string {
	return s.status.Err().Error()
}

func (s statusError) Problem(translator ut.Translator) iris.Problem {
	problem := iris.NewProblem()

	switch s.status.Code() {
	case codes.Unimplemented:
		problem.Status(iris.StatusNotImplemented)
	case codes.InvalidArgument:
		problem.Status(iris.StatusBadRequest)
	case codes.Unauthenticated:
		problem.Status(iris.StatusUnauthorized)
	case codes.FailedPrecondition:
		problem.Status(iris.StatusUnprocessableEntity)
	case codes.NotFound:
		problem.Status(iris.StatusNotFound)
	case codes.Aborted:
		problem.Status(iris.StatusConflict)
	default:
		problem.Status(iris.StatusInternalServerError)
	}

	problem.Detail(s.status.Message())

	for _, detail := range s.status.Details() {
		if badRequest, ok := detail.(*errdetails.BadRequest); ok {
			messages := make([]string, 0)
			for _, fieldViolation := range badRequest.GetFieldViolations() {
				messages = append(messages, fieldViolation.GetDescription())
			}

			problem.Key("errors", messages)
		}
	}

	return problem
}
