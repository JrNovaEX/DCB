// Package validator wraps go-playground/validator with DCB-specific rules
// and formats validation errors into user-friendly messages.
package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator is a configured instance of go-playground/validator.
type Validator struct {
	v *validator.Validate
}

// dockerImageRe is a permissive check for Docker image references.
var dockerImageRe = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9._\-/:@]+)?$`)

// New creates and returns a Validator with DCB-specific rules registered.
func New() (*Validator, error) {
	v := validator.New()

	if err := v.RegisterValidation("docker_image", validateDockerImage); err != nil {
		return nil, fmt.Errorf("registering docker_image validator: %w", err)
	}

	if err := v.RegisterValidation("port_range", validatePortRange); err != nil {
		return nil, fmt.Errorf("registering port_range validator: %w", err)
	}

	return &Validator{v: v}, nil
}

// Struct validates a struct using its validate tags.
// Returns a formatted error string listing all failures.
func (val *Validator) Struct(s interface{}) error {
	err := val.v.Struct(s)
	if err == nil {
		return nil
	}

	var errs validator.ValidationErrors
	if ok := isValidationErrors(err, &errs); !ok {
		return err
	}

	messages := make([]string, 0, len(errs))
	for _, fe := range errs {
		messages = append(messages, formatFieldError(fe))
	}
	return fmt.Errorf("validation failed:\n  %s", strings.Join(messages, "\n  "))
}

// validateDockerImage checks that a string looks like a valid Docker image reference.
func validateDockerImage(fl validator.FieldLevel) bool {
	return dockerImageRe.MatchString(fl.Field().String())
}

// validatePortRange checks that an int is in [1, 65535].
func validatePortRange(fl validator.FieldLevel) bool {
	p := fl.Field().Int()
	return p >= 1 && p <= 65535
}

// formatFieldError converts a single validation error into a readable message.
func formatFieldError(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("field %q is required", fe.Field())
	case "min":
		return fmt.Sprintf("field %q must have at least %s items", fe.Field(), fe.Param())
	case "docker_image":
		return fmt.Sprintf("field %q contains an invalid Docker image reference: %q", fe.Field(), fe.Value())
	case "port_range":
		return fmt.Sprintf("field %q must be between 1 and 65535, got %v", fe.Field(), fe.Value())
	default:
		return fmt.Sprintf("field %q failed %q validation", fe.Field(), fe.Tag())
	}
}

// isValidationErrors attempts to unwrap err as validator.ValidationErrors.
func isValidationErrors(err error, target *validator.ValidationErrors) bool {
	if ve, ok := err.(validator.ValidationErrors); ok {
		*target = ve
		return true
	}
	return false
}
