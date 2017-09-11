//go:generate mockgen -destination mock_multipartfile_test.go -package params_test mime/multipart File

package params_test

import (
	"errors"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	params "github.com/Nivl/go-params"
	"github.com/Nivl/go-params/formfile"
	"github.com/Nivl/go-params/formfile/mockformfile"
	"github.com/Nivl/go-params/formfile/testformfile"
	"github.com/Nivl/go-params/perror"
	"github.com/Nivl/go-types/date"
	"github.com/Nivl/go-types/filetype"
	"github.com/Nivl/go-types/ptrs"
)

func TestSetValue(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		t.Parallel()
		t.Run("regular", subTestsSetValueIntRegular)
		t.Run("pointer", subTestsSetValueIntPointer)
	})

	t.Run("string", func(t *testing.T) {
		t.Parallel()
		t.Run("regular", subTestsSetValueStringRegular)
		t.Run("pointer", subTestsSetValueStringPointer)
	})

	t.Run("bool", func(t *testing.T) {
		t.Parallel()
		t.Run("regular", subTestsSetValueBoolRegular)
		t.Run("pointer", subTestsSetValueStringPointer)
	})

	t.Run("scannable struct", subTestsSetValueScannableStruct)
}

func TestSetFile(t *testing.T) {
	t.Run("any type of files, required", subTestSetFileParamRequired)
	t.Run("only valid images", subTestSetFileParamValidImage)
	t.Run("ignore", subTestSetFileIgnore)
	t.Run("no name", subTestSetFileNoName)
	t.Run("wrong struct", subTestSetFileWrongStruct)
	t.Run("formFile returned an unknown error", subTestSetFileFormFileFail)
}

func subTestSetFileFormFileFail(t *testing.T) {
	t.Parallel()

	type strct struct {
		File string `from:"file"`
	}

	// Init the mocks
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Set the expectations
	fileHolder := mockformfile.NewMockFileHolder(mockCtrl)
	fileHolder.EXPECT().FormFile("File").Return(nil, nil, errors.New("unexpected error"))

	// Call the function to test
	paramList := reflect.ValueOf(&strct{}).Elem()
	p := newParamFromStructValue(&paramList, 0)
	err := p.SetFile(fileHolder)
	assert.Error(t, err, "Expected SetFile not to return an error")
}

func subTestSetFileWrongStruct(t *testing.T) {
	t.Parallel()

	type strct struct {
		File string `from:"file"`
	}

	// Init the mocks
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// create the multipart data
	cwd, _ := os.Getwd()
	imageHeader, imageFile := testformfile.NewMultipartData(t, cwd, "black_pixel.png")
	defer imageFile.Close()

	// Set the expectations
	fileHolder := mockformfile.NewMockFileHolder(mockCtrl)
	fileHolder.EXPECT().FormFile("File").Return(imageFile, imageHeader, nil)

	// Call the function to test
	paramList := reflect.ValueOf(&strct{}).Elem()
	p := newParamFromStructValue(&paramList, 0)
	err := p.SetFile(fileHolder)
	assert.Error(t, err, "Expected SetFile not to return an error")
	assert.True(t, strings.Contains(err.Error(), "the only accepted type for a file is"), "SetFile() failed with an unexpected error")
}

func subTestSetFileNoName(t *testing.T) {
	t.Parallel()

	type strct struct {
		File *formfile.FormFile `from:"file"`
	}

	// Init the mocks
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// create the multipart data
	cwd, _ := os.Getwd()
	imageHeader, imageFile := testformfile.NewMultipartData(t, cwd, "black_pixel.png")
	defer imageFile.Close()

	// Set the expectations
	fileHolder := mockformfile.NewMockFileHolder(mockCtrl)
	fileHolder.EXPECT().FormFile("File").Return(imageFile, imageHeader, nil)

	// Call the function to test
	paramList := reflect.ValueOf(&strct{}).Elem()
	p := newParamFromStructValue(&paramList, 0)
	err := p.SetFile(fileHolder)
	assert.NoError(t, err, "Expected SetFile not to return an error")
}

func subTestSetFileIgnore(t *testing.T) {
	t.Parallel()

	type strct struct {
		File *formfile.FormFile `from:"file" json:"-" params:"image"`
	}

	// init the mocks
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Set the expectations
	fileHolder := mockformfile.NewMockFileHolder(mockCtrl)

	// Call the function to test
	paramList := reflect.ValueOf(&strct{}).Elem()
	p := newParamFromStructValue(&paramList, 0)
	err := p.SetFile(fileHolder)
	assert.NoError(t, err, "SetFile() shuld not have fail")
}

func subTestSetFileParamValidImage(t *testing.T) {
	t.Parallel()

	type strct struct {
		File *formfile.FormFile `from:"file" json:"file" params:"image"`
	}

	testCases := []struct {
		description   string
		s             strct
		filename      string
		expectedMime  string
		expectedError error
	}{
		{"Valid PNG", strct{}, "black_pixel.png", "image/png", nil},
		{
			"invalid PNG",
			strct{},
			"invalid_content.png",
			"",
			perror.New("file", params.ErrMsgInvalidImage),
		},
		{
			"Not an image",
			strct{},
			"LICENSE",
			"",
			perror.New("file", filetype.ErrMsgUnsuportedImageFormat),
		},
		{
			"nil pointer should work as the image is not required",
			strct{},
			"",
			"",
			nil,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			// init the mocks
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// Set the expectations
			fileHolder := mockformfile.NewMockFileHolder(mockCtrl)
			formfile := fileHolder.EXPECT().FormFile("file")

			// if tc.filename is empty then we send no file
			if tc.filename == "" {
				formfile.Return(nil, nil, http.ErrMissingFile)
			} else {
				// create the multipart data
				cwd, _ := os.Getwd()
				imageHeader, imageFile := testformfile.NewMultipartData(t, cwd, tc.filename)
				defer imageFile.Close()

				formfile.Return(imageFile, imageHeader, nil)
			}

			// Call the function to test
			paramList := reflect.ValueOf(&tc.s).Elem()
			p := newParamFromStructValue(&paramList, 0)
			err := p.SetFile(fileHolder)

			// assert
			if tc.expectedError != nil {
				assert.Error(t, err, "Expected SetFile to return an error")
				assert.Equal(t, tc.expectedError, err, "Wrong error returned")
			} else {
				assert.NoError(t, err, "Expected SetFile not to return an error")

				if tc.filename != "" {
					assert.Equal(t, tc.expectedMime, tc.s.File.Mime, "Wrong mime type")
				}
			}
		})
	}
}

func subTestSetFileParamRequired(t *testing.T) {
	t.Parallel()

	type strct struct {
		File *formfile.FormFile `from:"file" json:"file" params:"required"`
	}

	testCases := []struct {
		description string
		s           strct
		sendNil     bool
		shouldFail  bool
	}{
		{"Nil pointer should fail", strct{}, true, true},
		{"Valid value should work", strct{}, false, false},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			// init the mocks
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// create the multipart data
			cwd, _ := os.Getwd()
			licenseHeader, licenseFile := testformfile.NewMultipartData(t, cwd, "LICENSE")
			defer licenseFile.Close()

			// Expectations
			fileHolder := mockformfile.NewMockFileHolder(mockCtrl)
			onFormFile := fileHolder.EXPECT().FormFile("file")

			if tc.sendNil {
				onFormFile.Return(nil, nil, http.ErrMissingFile)
			} else {
				onFormFile.Return(licenseFile, licenseHeader, nil)
			}

			paramList := reflect.ValueOf(&tc.s).Elem()
			p := newParamFromStructValue(&paramList, 0)

			err := p.SetFile(fileHolder)

			if tc.sendNil {
				assert.Error(t, err, "Expected SetFile to return an error")

			} else {
				assert.NoError(t, err, "Expected SetFile not to return an error")

				if assert.NotNil(t, tc.s.File, "Expected File NOT to be nil") {
					assert.NotNil(t, tc.s.File.File, "Expected File.File NOT to be nil")
					assert.NotNil(t, tc.s.File.Header, "Expected File.Header NOT to be nil")
					assert.Equal(t, licenseHeader.Filename, tc.s.File.Header.Filename)
				}
			}
		})
	}
}

func subTestsSetValueIntPointer(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description   string // optional, will use tag if empty
		tag           string
		source        url.Values
		expectedValue *int
	}{
		{
			"valid value, not tag",
			`json:"int"`,
			url.Values{"int": []string{"20"}},
			ptrs.NewInt(20),
		},
		{
			"not provided",
			`json:"int"`,
			url.Values{},
			nil,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			s := struct {
				Value *int
			}{}

			paramList := reflect.ValueOf(&s).Elem()
			p := newParamFromStructValue(&paramList, 0)
			tag := reflect.StructTag(tc.tag)
			p.Tags = &tag

			err := p.SetValue(tc.source)
			assert.NoError(t, err, "SetValue() should not have fail")

			if tc.expectedValue == nil {
				assert.Nil(t, s.Value, "SetValue() should not have set any value")
			} else {
				assert.Equal(t, *tc.expectedValue, *s.Value, "SetValue() did not set the expected value")
			}
		})
	}
}

func subTestsSetValueIntRegular(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description   string // optional, will use tag if empty
		tag           string
		source        url.Values
		expectedValue int
		expectedError error
	}{
		{
			"20 should work",
			`json:"int"`,
			url.Values{"int": []string{"20"}},
			20, nil,
		},
		{
			"ignored should work",
			`json:"-"`,
			url.Values{"value": []string{"20"}},
			0, nil,
		},
		{
			"20 with no name should work",
			``,
			url.Values{"Value": []string{"20"}},
			20, nil,
		},
		{
			"using default value should work",
			`json:"int" default:"42"`,
			url.Values{},
			42, nil,
		},
		{
			"-1 should work",
			`json:"int"`,
			url.Values{"int": []string{"-1"}},
			-1, nil,
		},
		{
			"not-an-int should fail",
			`json:"int"`,
			url.Values{"int": []string{"not-an-int"}},
			0,
			perror.New("int", params.ErrMsgInvalidInteger),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			s := struct {
				Value int
			}{}

			paramList := reflect.ValueOf(&s).Elem()
			p := newParamFromStructValue(&paramList, 0)
			tag := reflect.StructTag(tc.tag)
			p.Tags = &tag

			err := p.SetValue(tc.source)
			if tc.expectedError != nil {
				assert.Error(t, err, "SetValue() should have fail")
				assert.Equal(t, tc.expectedError, err, "SetValue() did not return the expected error")
			} else {
				assert.NoError(t, err, "SetValue() should not have fail")
				assert.Equal(t, tc.expectedValue, s.Value, "SetValue() did not set the expected value")
			}
		})
	}
}

func subTestsSetValueStringRegular(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description   string // optional, will use tag if empty
		tag           string
		source        url.Values
		expectedValue string
		expectedError error
	}{
		{
			"valid value should work",
			`json:"string"`,
			url.Values{"string": []string{"val"}},
			"val", nil,
		},
		{
			"default value should work",
			`json:"string" default:"default"`,
			url.Values{},
			"default", nil,
		},
		{
			"invalid uuid should fail with the uuid param",
			`json:"string" params:"uuid"`,
			url.Values{"string": []string{"no-a-uuid"}},
			"",
			perror.New("string", params.ErrMsgInvalidUUID),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			s := struct {
				Value string
			}{}

			paramList := reflect.ValueOf(&s).Elem()
			p := newParamFromStructValue(&paramList, 0)
			tag := reflect.StructTag(tc.tag)
			p.Tags = &tag

			err := p.SetValue(tc.source)
			if tc.expectedError != nil {
				assert.Error(t, err, "SetValue() should have fail")
				assert.Equal(t, tc.expectedError, err, "SetValue() did not return the expected error")
			} else {
				assert.NoError(t, err, "SetValue() should not have fail")
				assert.Equal(t, tc.expectedValue, s.Value, "SetValue() did not set the expected value")
			}
		})
	}
}

func subTestsSetValueStringPointer(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description   string // optional, will use tag if empty
		tag           string
		source        url.Values
		expectedValue *string
	}{
		{
			"valid value, not tag",
			`json:"string"`,
			url.Values{"string": []string{"val"}},
			ptrs.NewString("val"),
		},
		{
			"not provided",
			`json:"string"`,
			url.Values{},
			nil,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			s := struct {
				Value *string `json:"string"`
			}{}

			paramList := reflect.ValueOf(&s).Elem()
			p := newParamFromStructValue(&paramList, 0)
			tag := reflect.StructTag(tc.tag)
			p.Tags = &tag

			err := p.SetValue(tc.source)
			assert.NoError(t, err, "SetValue() should not have fail")

			if tc.expectedValue == nil {
				assert.Nil(t, s.Value, "SetValue() should not have set any value")
			} else {
				assert.Equal(t, *tc.expectedValue, *s.Value, "SetValue() did not set the expected value")
			}
		})
	}
}

func subTestsSetValueBoolRegular(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description   string // optional, will use tag if empty
		tag           string
		source        url.Values
		expectedValue bool
		expectedError error
	}{
		{
			"true should work",
			`json:"bool"`,
			url.Values{"bool": []string{"true"}},
			true, nil,
		},
		{
			"false should work",
			`json:"bool"`,
			url.Values{"bool": []string{"false"}},
			false, nil,
		},
		{
			"1 should work",
			`json:"bool"`,
			url.Values{"bool": []string{"1"}},
			true, nil,
		},
		{
			"0 should work",
			`json:"bool"`,
			url.Values{"bool": []string{"0"}},
			false, nil,
		},
		{
			"not value provided, using default value",
			`json:"bool" default:"true"`,
			url.Values{},
			true, nil,
		},
		{
			"invalid value should fail",
			`json:"bool"`,
			url.Values{"bool": []string{"not-a-bool"}},
			false,
			perror.New("bool", params.ErrMsgInvalidBoolean),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			s := struct {
				Value bool
			}{}

			paramList := reflect.ValueOf(&s).Elem()
			p := newParamFromStructValue(&paramList, 0)
			tag := reflect.StructTag(tc.tag)
			p.Tags = &tag

			err := p.SetValue(tc.source)
			if tc.expectedError != nil {
				assert.Error(t, err, "SetValue() should have fail")
				assert.Equal(t, tc.expectedError, err, "SetValue() did not return the expected error")
			} else {
				assert.NoError(t, err, "SetValue() should not have fail")
				assert.Equal(t, tc.expectedValue, s.Value, "SetValue() did not set the expected value")
			}
		})
	}
}

func subTestsSetValueBoolPointer(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description   string // optional, will use tag if empty
		tag           string
		source        url.Values
		expectedValue *bool
	}{
		{
			"true should work",
			`json:"bool"`,
			url.Values{"bool": []string{"true"}},
			ptrs.NewBool(true),
		},
		{
			"not provided should work",
			`json:"bool"`,
			url.Values{},
			nil,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			s := struct {
				Value *bool
			}{}

			paramList := reflect.ValueOf(&s).Elem()
			p := newParamFromStructValue(&paramList, 0)
			tag := reflect.StructTag(tc.tag)
			p.Tags = &tag

			err := p.SetValue(tc.source)
			assert.NoError(t, err, "SetValue() should not have fail")

			if tc.expectedValue == nil {
				assert.Nil(t, s.Value, "SetValue() should not have set any value")
			} else {
				assert.Equal(t, *tc.expectedValue, *s.Value, "SetValue() did not set the expected value")
			}
		})
	}
}

func subTestsSetValueScannableStruct(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description         string // optional, will use tag if empty
		tag                 string
		source              url.Values
		expectedStringValue string
		expectedError       error
	}{
		{
			"valid date should work",
			`json:"date"`,
			url.Values{"date": []string{"2017-09-10"}},
			"2017-09-10", nil,
		},
		{
			"invalid date should fail",
			`json:"date"`,
			url.Values{"date": []string{"not-a-date"}},
			"", perror.New("date", date.ErrMsgInvalidFormat),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			s := struct {
				Value *date.Date
			}{}

			paramList := reflect.ValueOf(&s).Elem()
			p := newParamFromStructValue(&paramList, 0)
			tag := reflect.StructTag(tc.tag)
			p.Tags = &tag

			err := p.SetValue(tc.source)
			if tc.expectedError != nil {
				assert.Error(t, err, "SetValue() should have fail")
				assert.Equal(t, tc.expectedError, err, "SetValue() did not return the expected error")
			} else {
				assert.NoError(t, err, "SetValue() should not have fail")
				assert.Equal(t, tc.expectedStringValue, s.Value.String(), "SetValue() did not set the expected value")
			}
		})
	}
}

// newParamFromStructValue creates a param using a struct value
func newParamFromStructValue(paramList *reflect.Value, paramPos int) *params.Param {
	value := paramList.Field(paramPos)
	info := paramList.Type().Field(paramPos)
	tags := info.Tag

	return &params.Param{
		Value: &value,
		Info:  &info,
		Tags:  &tags,
	}
}
