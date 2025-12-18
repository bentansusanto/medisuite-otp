package error

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	errValid "medisuite-api/constants/errors"

	"github.com/go-playground/validator/v10"
)

// ValidationError represents a single validation error.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ErrValidation is a map to hold custom validation error messages.
var ErrValidation = map[string]string{}

// ErrValidationResponse converts validation errors into a structured response.
func ErrValidationResponse(err error) (validationResponse []ValidationError) {
	var fieldErrors validator.ValidationErrors

	// Check if the error is of type ValidationErrors
	if errors.As(err, &fieldErrors) {
		for _, err := range fieldErrors {
			var message string
			switch err.Tag() {
			case "required":
				message = fmt.Sprintf(errValid.RequiredErrorMsg, err.Field())
			case "email":
				message = fmt.Sprintf(errValid.EmailErrorMsg, err.Field())
			case "oneof":
				message = fmt.Sprintf(errValid.OneOfErrorMsg, err.Field(), err.Param())
			default:
				// Check for custom error messages
				if errValidator, ok := ErrValidation[err.Tag()]; ok {
					count := strings.Count(errValidator, "%s")
					if count == 1 {
						message = fmt.Sprintf(errValidator, err.Field())
					} else {
						message = fmt.Sprintf(errValidator, err.Field(), err.Param())
					}
				} else {
					message = fmt.Sprintf("Something went wrong on %s; %s", err.Field(), err.Tag())
				}
			}
			validationResponse = append(validationResponse, ValidationError{
				Field:   err.Field(),
				Message: message,
			})
		}
	} else {
		// Log unexpected errors
		slog.Error("Validation error")
		validationResponse = append(validationResponse, ValidationError{
			Field:   "general",
			Message: errValid.GeneralError,
		})
	}

	return validationResponse
}

// WrapError logs the error and returns it.
func WrapError(err error) error {
	slog.Error("error validation")
	return err
}
