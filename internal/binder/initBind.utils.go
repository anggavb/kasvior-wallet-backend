package binder

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func InitValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "" {
				name = strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
			}
			if name == "-" {
				return ""
			}
			return name
		})

		_ = v.RegisterValidation("numeric", func(fl validator.FieldLevel) bool {
			if fl.Field().Kind() != reflect.String {
				return true
			}

			value := fl.Field().String()
			if value == "" {
				return false
			}

			for _, char := range value {
				if char < '0' || char > '9' {
					return false
				}
			}
			return true
		})
	}
}
