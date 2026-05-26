package binder

import (
	"log"
	"mime/multipart"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Initialize Binder Package
func init() {
	initValidate()
}

func initValidate() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Register Tag Name Func
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
		log.Println("RegisterTagNameFunc - get json/form struct tag")

		// Register Custom Validation for numeric tag
		v.RegisterValidation("numeric", func(fl validator.FieldLevel) bool {
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
		log.Println("RegisterValidation - Add custom validation for 'numeric' tag")

		// Register Custom Validator for image_check tag
		v.RegisterValidation("image_check", func(fl validator.FieldLevel) bool {
			file, ok := fl.Field().Interface().(*multipart.FileHeader)
			if !ok {
				return false
			}

			// Example: Max size 2MB
			if file.Size > 2*1024*1024 {
				return false
			}

			// Example: Allowed types
			allowedTypes := map[string]bool{
				"image/jpeg": true,
				"image/png":  true,
				"image/bmp":  true,
				"image/heic": true,
			}

			return allowedTypes[file.Header.Get("Content-Type")]
		})
		log.Println("RegisterValidation - Add custom validation for 'image_check' tag")
	}

}
