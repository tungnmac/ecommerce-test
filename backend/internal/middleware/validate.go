package middleware

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Custom validation function for strong password
var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterValidation("password", ValidatePassword)
}

func ValidatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 || len(password) > 32 {
		return false
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, ch := range password {
		switch {
		case 'A' <= ch && ch <= 'Z':
			hasUpper = true
		case 'a' <= ch && ch <= 'z':
			hasLower = true
		case '0' <= ch && ch <= '9':
			hasNumber = true
		case ch >= 33 && ch <= 47 || ch >= 58 && ch <= 64 || ch >= 91 && ch <= 96 || ch >= 123:
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// Reusable validator function
func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	errors := err.(validator.ValidationErrors)
	msg := ""
	for _, e := range errors {
		msg += fmt.Sprintf("Field '%s' failed on '%s' validation; ", e.Field(), e.Tag())
	}
	return fmt.Errorf(msg)
}
