package app

import (
	"github.com/go-playground/validator/v10"
	"regexp"
	"strings"
	"unicode"
)

func NewValidator() *validator.Validate {
	validate := validator.New()
	_ = validate.RegisterValidation("amikom", validateAmikomEmailOrUsername)
	_ = validate.RegisterValidation("containsany", containsAny)
	_ = validate.RegisterValidation("containslowercase", containsLowercase)
	_ = validate.RegisterValidation("containsuppercase", containsUppercase)
	_ = validate.RegisterValidation("containsnumeric", containsNumeric)

	return validate
}

func validateAmikomEmailOrUsername(field validator.FieldLevel) bool {
	emailOrUsername := field.Field().Interface().(string)
	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	isEmail, _ := regexp.MatchString(emailRegex, emailOrUsername)
	return isEmail && (strings.HasSuffix(emailOrUsername, "@students.amikom.ac.id") || strings.HasSuffix(emailOrUsername, "@amikom.ac.id")) || !isEmail
}

func containsAny(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	specialChars := "!@#$%^&*()_+~`"
	for _, char := range specialChars {
		if strings.ContainsRune(field, char) {
			return true
		}
	}
	return false
}

func containsLowercase(fl validator.FieldLevel) bool {
	return strings.ToLower(fl.Field().String()) != fl.Field().String()
}

func containsUppercase(fl validator.FieldLevel) bool {
	return strings.ToUpper(fl.Field().String()) != fl.Field().String()
}

func containsNumeric(fl validator.FieldLevel) bool {
	return strings.IndexFunc(fl.Field().String(), unicode.IsNumber) != -1
}
