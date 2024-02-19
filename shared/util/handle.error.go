package util

import (
	"github.com/go-playground/validator/v10"
)

type ApiErrorValidator struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func HandleErrorValidator(ve validator.ValidationErrors) []ApiErrorValidator {
	out := make([]ApiErrorValidator, len(ve))
	for i, fe := range ve {
		out[i] = ApiErrorValidator{fe.Field(), messageForTag(fe.Tag())}
	}
	return out
}

func messageForTag(tag string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "containsany":
		return "Value must contain at least one special character"
	case "containslowercase":
		return "Value must contain at least one lowercase character"
	case "containsuppercase":
		return "Value must contain at least one uppercase character"
	case "containsnumeric":
		return "Value must contain at least one number"
	case "eqfield":
		return "Value must be equal to the other field"
	case "min":
		return "Value must be at least 8 characters long"
	}
	return ""
}
