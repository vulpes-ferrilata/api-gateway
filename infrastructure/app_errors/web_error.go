package app_errors

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/kataras/iris/v12"
)

type WebError interface {
	Problem(translator ut.Translator) iris.Problem
}
