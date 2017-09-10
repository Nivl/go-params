package perror

import (
	"errors"
)

// Error is an interface represeting an error attached to a field name
type Error interface {
	error
	Field() string
}

// PError is an implementation of Error
type PError struct {
	error
	ErrorField string
}

// Field returns the field name attached to the error
func (err *PError) Field() string {
	return err.ErrorField
}

// New creates a new error using a field name and an error message
func New(field, errMsg string) *PError {
	return &PError{
		error:      errors.New(errMsg),
		ErrorField: field,
	}
}
