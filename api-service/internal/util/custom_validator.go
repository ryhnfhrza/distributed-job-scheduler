package util

import (
	"time"

	"github.com/go-playground/validator/v10"
)

func NewValidator(validate *validator.Validate) *validator.Validate {

	validate.RegisterValidation("future", func(fl validator.FieldLevel) bool {
		date, ok := fl.Field().Interface().(time.Time)
		if ok {
			return date.After(time.Now())
		}
		return false
	})

	return validate
}
