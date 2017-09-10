package params

const (
	// ErrMsgMissingParameter represents the error message corresponding to
	// a missing param
	ErrMsgMissingParameter = "parameter missing"

	// ErrMsgEmptyParameter represents the error message corresponding to
	// a missing param
	ErrMsgEmptyParameter = "parameter can be omitted but not empty"

	// ErrMsgInvalidUUID represents the error message corresponding to
	// an invalid UUID
	ErrMsgInvalidUUID = "not a valid uuid"

	// ErrMsgInvalidURL represents the error message corresponding to
	// an invalid URL
	ErrMsgInvalidURL = "not a valid url"

	// ErrMsgInvalidEmail represents the error message corresponding to
	// an invalid Email address
	ErrMsgInvalidEmail = "not a valid email"

	// ErrMsgInvalidImage represents the error message corresponding to
	// an invalid image
	ErrMsgInvalidImage = "not a valid image"

	// ErrMsgMaxLen represents the error message corresponding to
	// a field that exceed the maximum number of char
	ErrMsgMaxLen = "too many chars"

	// ErrMsgEnum represents the error message corresponding to
	// a field that doesn't contain a value set in an enum
	ErrMsgEnum = "not a valid value"

	// ErrMsgInvalidBoolean represents the error message corresponding to
	// an invalid boolean
	ErrMsgInvalidBoolean = "invalid boolean"

	// ErrMsgInvalidInteger represents the error message corresponding to
	// an invalid integer
	ErrMsgInvalidInteger = "invalid integer"
)
