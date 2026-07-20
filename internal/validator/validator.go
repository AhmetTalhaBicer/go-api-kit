package validator

import (
	"fmt"
	"strings"
)

// ValidationError holds one or more field-level validation errors.
type ValidationError struct {
	Fields map[string]string
}

func (e *ValidationError) Error() string {
	msgs := make([]string, 0, len(e.Fields))
	for field, msg := range e.Fields {
		msgs = append(msgs, fmt.Sprintf("%s: %s", field, msg))
	}
	return strings.Join(msgs, ", ")
}

// Validator is a fluent helper for validating struct fields.
type Validator struct {
	errors map[string]string
}

// New creates a new Validator instance.
func New() *Validator {
	return &Validator{errors: make(map[string]string)}
}

// Required checks that a string field is not empty.
func (v *Validator) Required(field, value string) *Validator {
	if strings.TrimSpace(value) == "" {
		v.errors[field] = "this field is required"
	}
	return v
}

// MinLength checks that a string meets a minimum length.
func (v *Validator) MinLength(field, value string, min int) *Validator {
	if len(value) < min {
		v.errors[field] = fmt.Sprintf("must be at least %d characters", min)
	}
	return v
}

// MaxLength checks that a string does not exceed a maximum length.
func (v *Validator) MaxLength(field, value string, max int) *Validator {
	if len(value) > max {
		v.errors[field] = fmt.Sprintf("must be at most %d characters", max)
	}
	return v
}

// Valid returns true when there are no validation errors.
func (v *Validator) Valid() bool {
	return len(v.errors) == 0
}

// Errors returns the collected validation errors, or nil if validation passed.
func (v *Validator) Errors() *ValidationError {
	if v.Valid() {
		return nil
	}
	return &ValidationError{Fields: v.errors}
}
