package params

import (
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/Nivl/go-params/perror"
	"github.com/Nivl/go-types/filetype"
	"github.com/Nivl/go-types/slices"
	"github.com/Nivl/go-types/strngs"
)

// Options represent all the options for a field
type Options struct {
	// Ignore means the field should not been parsed
	// json:"-"
	Ignore bool

	// Name contains the name of the field in the payload
	// json:"my_field"
	Name string

	// Required means the request should fail with a Bad Request if the field is missing.
	// params:"required"
	Required bool

	// Trim means the field needs to be trimmed before being retrieved and checked
	// params:"trim"
	Trim bool

	// ValidateUUID means the field should contain a valid UUIDv4
	// params:"uuid"
	ValidateUUID bool

	// ValidateEmail means the field should contain a valid email
	// params:"email"
	ValidateEmail bool

	// ValidateURL means the field should contain a valid url
	// params:"url"
	ValidateURL bool

	// ValidateImage means the field should contain a valid image
	// params:"image"
	ValidateImage bool

	// NoEmpty means the field should not contain an empty value
	// The difference between `required` and `noempty` is that `require` does`
	// not accept nil pointer. `noempty` accepts nil pointer, but if a value is
	// provided it cannot be an empty string.
	// params:"noempty"
	NoEmpty bool

	// MaxLen represents the maximum length a param can have (under its string
	// form). Any invalid values (including 0) will be ignored
	// maxlen:"255"
	MaxLen int

	// AuthorizedValues represents the list of authorized value for this param
	// enum:"and,or"
	AuthorizedValues []string
}

// NewOptions returns a ParamOptions from a StructTag
func NewOptions(tags *reflect.StructTag) *Options {
	output := &Options{}

	// We use the json tag to get the field name
	jsonOpts := strings.Split(tags.Get("json"), ",")
	if len(jsonOpts) > 0 {
		if jsonOpts[0] == "-" {
			return &Options{Ignore: true}
		}

		output.Name = jsonOpts[0]
	}

	// We use the maxlen tag to get the max length of a the value
	maxlen := tags.Get("maxlen")
	if len(maxlen) > 0 {
		// we silently fail on errors
		output.MaxLen, _ = strconv.Atoi(maxlen)
	}

	// We use the enu, tag to get all the authorized value a param can have
	enum := tags.Get("enum")
	if len(enum) > 0 {
		// we silently fail on errors
		output.AuthorizedValues = strings.Split(enum, ",")
	}

	// We parse the params
	opts := strings.Split(tags.Get("params"), ",")
	nbOptions := len(opts)
	for i := 0; i < nbOptions; i++ {
		switch opts[i] {
		case "required":
			output.Required = true
		case "noempty":
			output.NoEmpty = true
		case "trim":
			output.Trim = true
		case "uuid":
			output.ValidateUUID = true
		case "email":
			output.ValidateEmail = true
		case "url":
			output.ValidateURL = true
		case "image":
			output.ValidateImage = true
		}
	}

	return output
}

// Validate checks the given value passes the options set
func (opts *Options) Validate(value string, wasProvided bool) error {
	if opts.MaxLen > 0 && len(value) > opts.MaxLen {
		return perror.New(opts.Name, ErrMsgMaxLen)
	}

	if value == "" && opts.Required {
		return perror.New(opts.Name, ErrMsgMissingParameter)
	}

	if value == "" && opts.NoEmpty && wasProvided {
		return perror.New(opts.Name, ErrMsgEmptyParameter)
	}

	if value != "" {
		if opts.ValidateUUID && !strngs.IsValidUUID(value) {
			return perror.New(opts.Name, ErrMsgInvalidUUID)
		}

		if opts.ValidateURL && !strngs.IsValidURL(value) {
			return perror.New(opts.Name, ErrMsgInvalidURL)
		}

		if opts.ValidateEmail && !strngs.IsValidEmail(value) {
			return perror.New(opts.Name, ErrMsgInvalidEmail)
		}

		if len(opts.AuthorizedValues) > 0 {
			found, _ := slices.InSlice(opts.AuthorizedValues, value)
			if !found {
				return perror.New(opts.Name, ErrMsgEnum)
			}
		}
	}

	return nil
}

// ValidateFileContent checks the given file passes the options set
func (opts *Options) ValidateFileContent(file io.ReadSeeker) (string, error) {
	// Just for security, but it shouldn't be necessary
	defer file.Seek(0, io.SeekStart)

	if !opts.ValidateImage {
		// We still get the mimetype
		mimeType, err := filetype.MimeType(file)
		if err != nil {
			return "", err
		}
		return mimeType, nil
	}

	valid, mime, err := filetype.IsImage(file)
	if err != nil {
		if err.Error() == filetype.ErrMsgUnsuportedImageFormat {
			return "", perror.New(opts.Name, err.Error())
		}
		return "", err
	}
	if !valid {
		return "", perror.New(opts.Name, ErrMsgInvalidImage)
	}
	// check "valid", and return an error if its not
	return mime, nil
}

// ApplyTransformations applies all the wanted transformations to the given value
func (opts *Options) ApplyTransformations(value string) string {
	if opts.Trim {
		value = strings.TrimSpace(value)
	}
	return value
}
