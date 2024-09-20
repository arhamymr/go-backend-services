package helpers

import (
	"regexp"

	"github.com/go-playground/validator"
)

type AppValidator struct {
	Validator *validator.Validate
}

func (cv *AppValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

func ValidatePassword(f validator.FieldLevel) bool {
	password := f.Field().String()
	var (
		hasMinLen  = len(password) >= 8
		hasUpper   = regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLower   = regexp.MustCompile(`[a-z]`).MatchString(password)
		hasNumber  = regexp.MustCompile(`[0-9]`).MatchString(password)
		hasSpecial = regexp.MustCompile(`[!@#~$%^&*()_+|<>?:{}]`).MatchString(password)
	)
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}
