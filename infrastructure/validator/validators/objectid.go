package validators

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RegisterObjectIDValidator(v *validator.Validate) error {
	if err := v.RegisterValidation("objectid", isObjectID); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func isObjectID(fl validator.FieldLevel) bool {
	field := fl.Field()

	switch field.Kind() {
	case reflect.String:
		return primitive.IsValidObjectID(field.String())
	default:
		return false
	}
}
