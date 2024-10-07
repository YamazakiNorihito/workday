package validation_error

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	errors map[string]string
}

func New(errors map[string]string) *ValidationError {
	return &ValidationError{errors: errors}
}

func (ve *ValidationError) Error() string {
	var errMessages []string
	for field, message := range ve.errors {
		errMessages = append(errMessages, fmt.Sprintf("%s: %s", field, message))
	}
	return fmt.Sprintf("Validation failed: %s", strings.Join(errMessages, ", "))
}

func (ve *ValidationError) Errors() map[string]string {
	return ve.errors
}
