package request

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func isValid[T any](payload T) error {
	err := validate.Struct(payload)
	if err != nil {
		var errorMessages []string
		for _, fieldErr := range err.(validator.ValidationErrors) {
			message := getCustomMessage(fieldErr)
			errorMessages = append(errorMessages, message)
		}
		return fmt.Errorf("validation errors: %s", errorMessages)
	}
	return nil
}

func getCustomMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is a required field", fe.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", fe.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", fe.Field(), fe.Param())
	default:
		return fmt.Sprintf("%s is invalid", fe.Field())
	}
}
