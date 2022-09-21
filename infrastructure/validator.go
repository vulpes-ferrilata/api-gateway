package infrastructure

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/pkg/errors"
	custom_en_translations "github.com/vulpes-ferrilata/api-gateway/infrastructure/validator/translations/en"
	custom_validators "github.com/vulpes-ferrilata/api-gateway/infrastructure/validator/validators"
)

var (
	ErrTranslatorNotFound = errors.New("translator not found")
)

func NewValidator(universalTranslator *ut.UniversalTranslator) (*validator.Validate, error) {
	validate := validator.New()

	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	if err := custom_validators.RegisterObjectIDValidator(validate); err != nil {
		return nil, errors.WithStack(err)
	}

	en := en.New()

	translator, ok := universalTranslator.GetTranslator(en.Locale())
	if !ok {
		return nil, errors.Wrap(ErrTranslatorNotFound, en.Locale())
	}

	if err := en_translations.RegisterDefaultTranslations(validate, translator); err != nil {
		return nil, errors.WithStack(err)
	}

	if err := custom_en_translations.RegisterDefaultTranslations(validate, translator); err != nil {
		return nil, errors.WithStack(err)
	}

	return validate, nil
}
