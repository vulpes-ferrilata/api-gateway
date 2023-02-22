package app_errors

import (
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
)

func NewRequestValidationError(validationErrors validator.ValidationErrors) AppError {
	return &requestValidationError{
		validationErrors: validationErrors,
	}
}

type requestValidationError struct {
	validationErrors validator.ValidationErrors
}

func (r requestValidationError) Error() string {
	builder := new(strings.Builder)

	builder.WriteString("the request contains invalid parameters")

	for _, fieldError := range r.validationErrors {
		builder.WriteString("\n")
		builder.WriteString(fieldError.Error())
	}

	return builder.String()
}

func (r requestValidationError) Problem(translator ut.Translator) iris.Problem {
	problem := iris.NewProblem()

	problem.Status(iris.StatusBadRequest)

	detail, err := translator.T("request-validation-error")
	if err != nil {
		detail = "the request contains invalid parameters"
	}
	problem.Detail(detail)

	messages := make([]string, 0)
	for _, fieldError := range r.validationErrors {
		messages = append(messages, fieldError.Translate(translator))
	}
	problem.Key("errors", messages)

	return problem
}
