package infrastructure

import (
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/vi"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	vi_translations "github.com/go-playground/validator/v10/translations/vi"
	"github.com/pkg/errors"
	custom_en_translations "github.com/vulpes-ferrilata/custom-validator/translations/en"
	custom_vi_translations "github.com/vulpes-ferrilata/custom-validator/translations/vi"
	custom_validators "github.com/vulpes-ferrilata/custom-validator/validators"
)

var (
	ErrTranslatorNotFound = errors.New("translator not found")
)

func NewValidator(universalTranslator *ut.UniversalTranslator) (*validator.Validate, error) {
	validate := validator.New()

	if err := custom_validators.RegisterObjectIDValidator(validate); err != nil {
		return nil, errors.WithStack(err)
	}

	en := en.New()
	enTranslator, ok := universalTranslator.GetTranslator(en.Locale())
	if !ok {
		return nil, errors.Wrap(ErrTranslatorNotFound, en.Locale())
	}

	if err := en_translations.RegisterDefaultTranslations(validate, enTranslator); err != nil {
		return nil, errors.WithStack(err)
	}

	if err := custom_en_translations.RegisterDefaultTranslations(validate, enTranslator); err != nil {
		return nil, errors.WithStack(err)
	}

	vi := vi.New()
	viTranslator, ok := universalTranslator.GetTranslator(vi.Locale())
	if !ok {
		return nil, errors.Wrap(ErrTranslatorNotFound, vi.Locale())
	}

	if err := vi_translations.RegisterDefaultTranslations(validate, viTranslator); err != nil {
		return nil, errors.WithStack(err)
	}

	if err := custom_vi_translations.RegisterDefaultTranslations(validate, viTranslator); err != nil {
		return nil, errors.WithStack(err)
	}

	return validate, nil
}
