package params

// CustomValidation is an interface used to implements custom validation on
// a structure
type CustomValidation interface {
	// IsValid checks if the provided params are valid.
	// returns a boolean, the name of the field that failed (if any),
	// and the error (if any)
	IsValid() (isValid bool, fieldFailing string, err error)
}
