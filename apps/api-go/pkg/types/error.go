package types

import (
	"fmt"
	"net/http"
)

// FieldValidationError represents a validation error
type FieldValidationError interface {
	error
	HTTPStatus() int
	Type() string
	GetLocation() string
	GetValue() any
	GetMessage() string
}

// ValidationError represents a validation error
type ValidationError struct {
	Message  string
	Field    string
	Location string
	Value    any
}

// Error implements error interface
func (v *ValidationError) Error() string {
	if v.Field != "" {
		return fmt.Sprintf("validation error for field %s: %s", v.Field, v.Message)
	}
	return fmt.Sprintf("validation error: %s", v.Message)
}

// HTTPStatus returns the HTTP status code for the error
func (v *ValidationError) HTTPStatus() int {
	return http.StatusUnprocessableEntity
}

// Type returns the type of the error
func (v *ValidationError) Type() string {
	return "VALIDATION_ERROR"
}

// GetLocation returns the location of the error
func (v *ValidationError) GetLocation() string {
	if v.Location != "" {
		return v.Location
	}
	if v.Field != "" {
		return fmt.Sprintf("body.%s", v.Field)
	}
	return ""
}

// GetValue returns the value that caused the error
func (v *ValidationError) GetValue() any {
	return v.Value
}

// GetMessage returns the error message
func (v *ValidationError) GetMessage() string {
	return v.Message
}

// NewValidationError creates a new validation error
func NewValidationError(message string, field string, location string, value any) *ValidationError {
	return &ValidationError{
		Message:  message,
		Field:    field,
		Location: location,
		Value:    value,
	}
}

// ServiceError represents a service error
type ServiceError interface {
	error
	HTTPStatus() int
	Type() string
}

// NotFoundError represents a not found error
type NotFoundError struct {
	Message string
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{
		Message: message,
	}
}

// Error implements error interface
func (n *NotFoundError) Error() string {
	return n.Message
}

// HTTPStatus returns the HTTP status code for the error
func (n *NotFoundError) HTTPStatus() int {
	return http.StatusNotFound
}

// Type returns the type of the error
func (n *NotFoundError) Type() string {
	return "NOT_FOUND_ERROR"
}
