//go:generate mockgen -destination mock_readseekder_test.go -package params_test io ReadSeeker

package params_test

import (
	"errors"
	"io"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/Nivl/go-types/ptrs"

	"github.com/Nivl/go-types/filetype"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	params "github.com/Nivl/go-params"
	"github.com/Nivl/go-params/perror"
)

func TestNewOptions(t *testing.T) {
	testCases := []struct {
		description string // optional, will use tag if empty
		tag         string
		expected    *params.Options
	}{
		{
			"Empty tag", ``, &params.Options{},
		},
		{
			"Set Ignore", `json:"-"`,
			&params.Options{
				Ignore: true,
			},
		},
		{
			"Set Name", `json:"var_name"`,
			&params.Options{
				Name: "var_name",
			},
		},
		{
			"Set MaxLen", `maxlen:"255"`,
			&params.Options{
				MaxLen: 255,
			},
		},
		{
			"Empty MaxLen should be ignored", `maxlen:""`,
			&params.Options{},
		},
		{
			"Set AuthorizedValues", `enum:"val1,val2,val3"`,
			&params.Options{
				AuthorizedValues: []string{"val1", "val2", "val3"},
			},
		},
		{
			"Set empty AuthorizedValues", `enum:""`,
			&params.Options{},
		},
		{
			"Set ValidateUUID", `params:"uuid"`,
			&params.Options{
				ValidateUUID: true,
			},
		},
		{
			"Set Required", `params:"required"`,
			&params.Options{
				Required: true,
			},
		},
		{
			"Set Required", `params:"required"`,
			&params.Options{
				Required: true,
			},
		},
		{
			"Set NoEmpty", `params:"noempty"`,
			&params.Options{
				NoEmpty: true,
			},
		},
		{
			"Set Trim", `params:"trim"`,
			&params.Options{
				Trim: true,
			},
		},
		{
			"Set ValidateEmail", `params:"email"`,
			&params.Options{
				ValidateEmail: true,
			},
		},
		{
			"Set ValidateURL", `params:"url"`,
			&params.Options{
				ValidateURL: true,
			},
		},
		{
			"Set ValidateSlug", `params:"slug"`,
			&params.Options{
				ValidateSlug: true,
			},
		},
		{
			"Set ValidateImage", `params:"image"`,
			&params.Options{
				ValidateImage: true,
			},
		},
		{
			"Set MaxInt", `max_int:"1"`,
			&params.Options{
				MaxInt: ptrs.NewInt(1),
			},
		},
		{
			"Set MinInt", `min_int:"-1"`,
			&params.Options{
				MinInt: ptrs.NewInt(-1),
			},
		},
		{
			"", `json:"my_var" params:"email,required" maxlen:"30"`,
			&params.Options{
				Name:          "my_var",
				ValidateEmail: true,
				Required:      true,
				MaxLen:        30,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		if tc.description == "" {
			tc.description = tc.tag
		}

		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			tag := reflect.StructTag(tc.tag)
			output, err := params.NewOptions(&tag)
			require.NoError(t, err, "NewOptions() should not have failed)")
			assert.Equal(t, tc.expected, output, "NewOptions()returned unexpected data")
		})
	}
}

func TestNewOptionsError(t *testing.T) {
	testCases := []struct {
		description string // optional, will use tag if empty
		tag         string
	}{
		{
			"Set MaxInt", `max_int:"nan"`,
		},
		{
			"Set MinInt", `min_int:"nan"`,
		},
		{
			"Set MaxLen", `maxlen:"nan"`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		if tc.description == "" {
			tc.description = tc.tag
		}

		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			tag := reflect.StructTag(tc.tag)
			_, err := params.NewOptions(&tag)
			require.Error(t, err, "NewOptions() should not have failed)")
		})
	}
}

func TestValidate(t *testing.T) {
	// sugars
	wasProvided := true

	testCases := []struct {
		description   string
		tag           string
		value         string
		wasProvided   bool
		expectedError error
	}{
		{
			"maxlen with valid data should work",
			`json:"" maxlen:"3"`,
			"val",
			wasProvided,
			nil,
		},
		{
			"maxlen with valid data should fail",
			`json:"field_name" maxlen:"3"`,
			"too many chars",
			wasProvided,
			perror.New("field_name", params.ErrMsgMaxLen),
		},
		{
			"required with valid data",
			`json:"field_name" params:"required"`,
			"valid data",
			wasProvided,
			nil,
		},
		{
			"required with invalid data",
			`json:"field_name" params:"required"`,
			"",
			wasProvided,
			perror.New("field_name", params.ErrMsgMissingParameter),
		},
		{
			"noempty with valid data",
			`json:"field_name" params:"noempty"`,
			"valid data",
			wasProvided,
			nil,
		},
		{
			"noempty with no provided data",
			`json:"field_name" params:"noempty"`,
			"",
			!wasProvided,
			nil,
		},
		{
			"noempty with invalid, provided data",
			`json:"field_name" params:"noempty"`,
			"",
			wasProvided,
			perror.New("field_name", params.ErrMsgEmptyParameter),
		},
		{
			"uuid with valid data",
			`json:"field_name" params:"uuid,required"`,
			"b3ca2cb7-422c-4467-a3ed-bce00a6a8216",
			wasProvided,
			nil,
		},
		{
			"uuid with invalid data",
			`json:"field_name" params:"uuid,required"`,
			"not-a-uuid",
			wasProvided,
			perror.New("field_name", params.ErrMsgInvalidUUID),
		},
		{
			"slug with valid data",
			`json:"field_name" params:"slug,required"`,
			"this-is-a-val1d-slug",
			wasProvided,
			nil,
		},
		{
			"slug with invalid data",
			`json:"field_name" params:"slug,required"`,
			"not a SLUG",
			wasProvided,
			perror.New("field_name", params.ErrMsgInvalidSlug),
		},
		{
			"slugOrUuid with valid data (slug)",
			`json:"field_name" params:"slugOrUuid"`,
			"this-is-a-val1d-slug",
			wasProvided,
			nil,
		},
		{
			"slugOrUuid with valid data (uuid)",
			`json:"field_name" params:"slugOrUuid"`,
			"b3ca2cb7-422c-4467-a3ed-bce00a6a8216",
			wasProvided,
			nil,
		},
		{
			"slugOrUuid with invalid data",
			`json:"field_name" params:"slugOrUuid"`,
			"not a SLUG or UUID",
			wasProvided,
			perror.New("field_name", params.ErrMsgInvalidSlugOrUUID),
		},
		{
			"url with valid data",
			`json:"field_name" params:"url,required"`,
			"https://google.com",
			wasProvided,
			nil,
		},
		{
			"url with invalid data",
			`json:"field_name" params:"url,required"`,
			"not-a-url",
			wasProvided,
			perror.New("field_name", params.ErrMsgInvalidURL),
		},
		{
			"email with valid data",
			`json:"field_name" params:"email,required"`,
			"hi@melvin.la",
			wasProvided,
			nil,
		},
		{
			"email with invalid data",
			`json:"field_name" params:"email,required"`,
			"hi.melvin.la",
			wasProvided,
			perror.New("field_name", params.ErrMsgInvalidEmail),
		},
		{
			"enum with valid data",
			`json:"field_name" enum:"val1,va2" params:"required"`,
			"val1",
			wasProvided,
			nil,
		},
		{
			"enum with invalid data",
			`json:"field_name" enum:"val1,va2" params:"required"`,
			"not a valid value",
			wasProvided,
			perror.New("field_name", params.ErrMsgEnum),
		},
		{
			"min_int with valid data should work",
			`json:"field_name" min_int:"-1"`,
			"-1",
			wasProvided,
			nil,
		},
		{
			"min_int with invalid type should fail",
			`json:"field_name" min_int:"-1"`,
			"nan",
			wasProvided,
			perror.New("field_name", params.ErrMsgInvalidInteger),
		},
		{
			"min_int with invalid data should fail",
			`json:"field_name" min_int:"-1"`,
			"-2",
			wasProvided,
			perror.New("field_name", params.ErrMsgIntegerTooSmall),
		},
		{
			"max_int with valid data should work",
			`json:"field_name" max_int:"1"`,
			"1",
			wasProvided,
			nil,
		},
		{
			"max_int with invalid type should fail",
			`json:"field_name" max_int:"1"`,
			"nan",
			wasProvided,
			perror.New("field_name", params.ErrMsgInvalidInteger),
		},
		{
			"max_int with invalid data should fail",
			`json:"field_name" max_int:"1"`,
			"2",
			wasProvided,
			perror.New("field_name", params.ErrMsgIntegerTooBig),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			tag := reflect.StructTag(tc.tag)
			opts, err := params.NewOptions(&tag)
			require.NoError(t, err, "NewOptions() should not have failed)")

			err = opts.Validate(tc.value, tc.wasProvided)
			if tc.expectedError != nil {
				assert.Error(t, err, "Validate() should have failed")
				assert.Equal(t, tc.expectedError, err, "Validate() returned an unexcpected error")
			} else {
				assert.NoError(t, err, "Validate() should not have failed")
			}
		})
	}
}

func TestApplyTransformations(t *testing.T) {
	testCases := []struct {
		description string
		tag         string
		input       string
		output      string
	}{
		{
			"trim a string containing only spaces",
			`params:"trim"`,
			"        ",
			"",
		},
		{
			"trim a string",
			`params:"trim"`,
			"    test    ",
			"test",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			tag := reflect.StructTag(tc.tag)
			opts, err := params.NewOptions(&tag)
			require.NoError(t, err, "NewOptions() should not have failed)")

			output := opts.ApplyTransformations(tc.input)
			assert.Equal(t, tc.output, output, "ApplyTransformations() returned an unexpected output")
		})
	}
}

func TestValidateFileContent(t *testing.T) {
	testCases := []struct {
		description      string
		tag              string
		filename         string
		expectedMime     string
		ExpectedErrorMsg string
	}{
		{
			"validate a valid image",
			`params:"image"`,
			`black_pixel.png`,
			"image/png",
			"",
		},
		{
			"don't validate a valid image",
			``,
			`black_pixel.png`,
			"image/png",
			"",
		},
		{
			"validate an invalid image",
			`params:"image"`,
			`invalid_content.png`,
			"",
			params.ErrMsgInvalidImage,
		},
		{
			"validate an invalid image without checking the type",
			``,
			`invalid_magic.png`,
			"application/octet-stream",
			"",
		},
		{
			"validate an invalid image and check the type",
			`params:"image"`,
			`invalid_magic.png`,
			"",
			filetype.ErrMsgUnsuportedImageFormat,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			tag := reflect.StructTag(tc.tag)
			opts, err := params.NewOptions(&tag)
			require.NoError(t, err, "NewOptions() should not have failed)")
			filepath := path.Join("fixtures", tc.filename)
			file, err := os.Open(filepath)
			require.NoError(t, err, "Open() should not have failed)")

			mime, err := opts.ValidateFileContent(file)
			if tc.ExpectedErrorMsg != "" {
				require.Error(t, err, "ValidateFileContent() should have failed")
				assert.Equal(t, tc.ExpectedErrorMsg, err.Error(), "ValidateFileContent() did not fail with the expected error")
			} else {
				require.NoError(t, err, "ValidateFileContent() should not have failed")
				assert.Equal(t, tc.expectedMime, mime, "ValidateFileContent() did not return the expected mime")
			}
		})
	}
}

func TestValidateFileContentMock(t *testing.T) {
	t.Run("fail reading mime type when no validating", func(t *testing.T) {
		t.Parallel()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		exectedError := errors.New("read failed")

		// Set the expectations
		reader := NewMockReadSeeker(mockCtrl)
		reader.EXPECT().Seek(int64(0), io.SeekStart).Return(int64(0), nil)
		reader.EXPECT().Read(gomock.Any()).Return(0, exectedError)

		opts := &params.Options{}
		_, err := opts.ValidateFileContent(reader)
		require.Error(t, err, "ValidateFileContent() should have failed")
		assert.Equal(t, exectedError, err, "ValidateFileContent() did not fail with the expected error")
	})

	t.Run("fail reading anythong when validating", func(t *testing.T) {
		t.Parallel()
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		exectedError := errors.New("read failed")

		// Set the expectations
		reader := NewMockReadSeeker(mockCtrl)
		reader.EXPECT().Seek(int64(0), io.SeekStart).Return(int64(0), nil)
		reader.EXPECT().Read(gomock.Any()).Return(0, exectedError)

		opts := &params.Options{ValidateImage: true}
		_, err := opts.ValidateFileContent(reader)
		require.Error(t, err, "ValidateFileContent() should have failed")
		assert.Equal(t, exectedError, err, "ValidateFileContent() did not fail with the expected error")
	})
}
