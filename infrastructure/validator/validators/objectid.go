package validators

import (
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
	return primitive.IsValidObjectID(fl.Field().String())
}
