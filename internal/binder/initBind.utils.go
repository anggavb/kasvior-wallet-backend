package binder

import (
	"log"
	"mime/multipart"
	"reflect"
	"strconv"
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
		v.RegisterValidation("image_max_size", func(fl validator.FieldLevel) bool {
			file, ok := fl.Field().Interface().(multipart.FileHeader)
			if !ok {
				log.Println("Invalid type for image_check validation")
				return false
			}

			// Example: Max size 2MB
			param := fl.Param()
			log.Println(param)
			maxSize, err := strconv.ParseInt(param, 10, 64)
			if err != nil {
				return false // Gagalkan validasi jika parameter tag bukan angka konkrit
			}

			if file.Size > maxSize {
				log.Println("Max size more than allowed")
				return false
			}

			return true
		})
		log.Println("RegisterValidation - Add custom validation for 'image_max_size' tag")

		v.RegisterValidation("image_type", func(fl validator.FieldLevel) bool {
			file, ok := fl.Field().Interface().(multipart.FileHeader)
			if !ok {
				log.Println("Invalid type for image_type validation")
				return false
			}

			// Example: Allowed types
			allowedTypes := map[string]bool{
				"image/jpeg": true,
				"image/png":  true,
				"image/bmp":  true,
				"image/heic": true,
				"image/webp": true,
			}

			return allowedTypes[file.Header.Get("Content-Type")]
		})
		log.Println("RegisterValidation - Add custom validation for 'image_type' tag")
	}
}
