package infrastructure

import (
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/vi"
	ut "github.com/go-playground/universal-translator"
	"github.com/pkg/errors"
)

func NewUniversalTranslator() (*ut.UniversalTranslator, error) {
	enLocale := en.New()
	viLocale := vi.New()

	universalTranslator := ut.New(enLocale, enLocale, viLocale)

	if err := universalTranslator.Import(ut.FormatJSON, "./locales"); err != nil {
		return nil, errors.WithStack(err)
	}

	return universalTranslator, nil
}
