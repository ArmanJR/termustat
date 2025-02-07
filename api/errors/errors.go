package errors

import (
	stdErrors "errors"
	"fmt"
)

// Base error types
type Error string

func (e Error) Error() string { return string(e) }

// Predefined errors
const (
	ErrNotFound   = Error("record not found")
	ErrConflict   = Error("record already exists")
	ErrInvalid    = Error("invalid input")
	ErrForbidden  = Error("forbidden")
	ErrBadRequest = Error("bad request")
)

// Custom error types
type NotFoundError struct {
	Entity string
	ID     string
	err    error
}
type NotValidError struct {
	Entity string
	err    error
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %s not found", e.Entity, e.ID)
}
func (e *NotFoundError) Unwrap() error {
	return e.err
}

func (e *NotValidError) Error() string {
	return fmt.Sprintf("%s is not a valid input", e.Entity)
}
func (e *NotValidError) Unwrap() error {
	return e.err
}

// Error constructors
func NewNotFoundError(entity, id string) error {
	return &NotFoundError{
		Entity: entity,
		ID:     id,
		err:    ErrNotFound,
	}
}
func NewValidationError(entity string) error {
	return &NotValidError{
		Entity: entity,
		err:    ErrInvalid,
	}
}

// Wrap wraps an error with additional context
func Wrap(err error, message string) error {
	return fmt.Errorf("%s: %w", message, err)
}

// Is reports whether any error in err's chain matches target
func Is(err, target error) bool {
	return stdErrors.Is(err, target)
}

// As finds the first error in err's chain that matches target
func As(err error, target any) bool {
	return stdErrors.As(err, target)
}

// New creates a new error with the given message
func New(text string) error {
	return stdErrors.New(text)
}
