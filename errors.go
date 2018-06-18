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

	// ErrMsgInvalidSlug represents the error message corresponding to
	// an invalid slug
	ErrMsgInvalidSlug = "not a valid slug"

	// ErrMsgInvalidSlugOrUUID represents the error message corresponding to
	// a field that is neither a slug or a UUIDv4
	ErrMsgInvalidSlugOrUUID = "not a valid slug"

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

	// ErrMsgIntegerTooBig represents the error message corresponding to
	// an integer being too big
	ErrMsgIntegerTooBig = "value too high"

	// ErrMsgIntegerTooSmall represents the error message corresponding to
	// an integer being too small
	ErrMsgIntegerTooSmall = "value too small"

	// ErrMsgEmptyFile represents the error message corresponding to
	// an empty file being sent
	ErrMsgEmptyFile = "file empty"

	// ErrMsgCorruptedFile represents the error message corresponding to
	// a corrupted file
	ErrMsgCorruptedFile = "file seems corrupted"

	// ErrMsgArrayTooBig represents the error message corresponding to
	// an array being too big
	ErrMsgArrayTooBig = "too many elements"

	// ErrMsgArrayTooSmall represents the error message corresponding to
	// an array being too small
	ErrMsgArrayTooSmall = "too few elements"

	// ErrMsgEmptyItem represents the error message corresponding to
	// an array containing an empty item
	ErrMsgEmptyItem = "array cannot contain empty items"
)
