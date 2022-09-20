package infrastructure

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/pkg/errors"
)

func NewUniversalTranslator() (*ut.UniversalTranslator, error) {
	en := en.New()

	universalTranslator := ut.New(en, en)

	if err := universalTranslator.Import(ut.FormatJSON, "./locales"); err != nil {
		return nil, errors.WithStack(err)
	}

	return universalTranslator, nil
}
