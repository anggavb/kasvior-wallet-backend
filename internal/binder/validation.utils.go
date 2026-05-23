package binder

import (
	"errors"
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
)

func FormatValidationError(err error) map[string]string {
	out := make(map[string]string)

	if validationError, ok := errors.AsType[validator.ValidationErrors](err); ok {
		for _, e := range validationError {
			field := e.Field()
			tag := e.Tag()
			param := e.Param()

			switch tag {
			case "required":
				out[field] = fmt.Sprintf("The %s field is required.", field)
			case "email":
				out[field] = fmt.Sprintf("The %s field is not a valid email address.", field)
			case "gte":
				out[field] = fmt.Sprintf("The %s field must be greater than or equal to %s.", field, param)
			case "lte":
				out[field] = fmt.Sprintf("The %s field must be less than or equal to %s.", field, param)
			case "min":
				out[field] = fmt.Sprintf("The %s field must be at least %s characters long.", field, param)
			case "max":
				out[field] = fmt.Sprintf("The %s field must be at most %s characters long.", field, param)
			case "len":
				out[field] = fmt.Sprintf("The %s field must be exactly %s characters long.", field, param)
			case "oneof":
				out[field] = fmt.Sprintf("The %s field must be one of the following: %s.", field, param)
			default:
				out[field] = fmt.Sprintf("The %s field does not meet the validation requirements for %s %s.", field, tag, param)
			}
		}
		return out
	}

	log.Printf("Unexpected error: %v", err)
	out["error"] = "An unexpected error occurred during validation."
	return out
}
