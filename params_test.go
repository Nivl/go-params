package params_test

import (
	"errors"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"testing"

	params "github.com/Nivl/go-params"
	"github.com/Nivl/go-params/formfile"
	"github.com/Nivl/go-params/formfile/mockformfile"
	"github.com/Nivl/go-params/formfile/testformfile"
	"github.com/Nivl/go-types/date"
	"github.com/Nivl/go-types/ptrs"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type StructWithValidator struct {
	String     string `from:"query" json:"string" default:"default value"`
	TrueToFail bool   `from:"query" json:"true_to_fail" default:"false"`
}

func (p *StructWithValidator) IsValid() (isValid bool, fieldName string, err error) {
	if !p.TrueToFail {
		return true, "", nil
	}
	return false, "true_to_fail", errors.New("cannot be set to true")
}

func TestParamsParse(t *testing.T) {
	t.Run("valid struct", subTestValidStruct)
	t.Run("invalid struct", subTestInvalidStruct)
	t.Run("field with no source", subTestFieldWithNoSource)
	t.Run("field with unexisting source", subTestFieldWithUnexistingSource)
	t.Run("field not exported", subTestFieldNotExported)
	t.Run("embedded struct", subTestEmbeddedStruct)
	t.Run("embedded struct with custon validation", subTestEmbeddedStructWithCustomValidation)
	t.Run("custom validation", subTestCustomValidation)
	t.Run("file handling", subTestFileUpload)
	t.Run("file handling", subTestFileUpload)
}

func TestParamsExtract(t *testing.T) {
	t.Run("extract", subTestExtraction)
	t.Run("nil value", subTestExtractNil)
}

func subTestValidStruct(t *testing.T) {
	t.Parallel()

	type strct struct {
		ID            string  `from:"url" json:"id" params:"uuid,required"`
		Number        int     `from:"query" json:"number"`
		IgnoredInt    int     `from:"form" json:"-" params:"required"`
		RequiredBool  bool    `from:"form" json:"required_bool" params:"required"`
		PointerBool   *bool   `from:"form" json:"pointer_bool"`
		PointerString *string `from:"form" json:"pointer_string" params:"trim"`
		Default       int     `from:"form" json:"default" default:"42"`
		Emum          int     `from:"form" json:"enum" enum:"21,42"`
		DefaultSlice  []int   `from:"form" json:"default_slice" default:"1,2,3"`
	}

	s := &strct{}
	p := params.New(s)

	urlSource := url.Values{}
	urlSource.Set("id", "1aa75114-6117-4908-b6ea-0d22ecdd4fc0")

	querySource := url.Values{}
	querySource.Set("number", "24")

	formSource := url.Values{}
	formSource.Set("IgnoredInt", "42")
	formSource.Set("required_bool", "true")
	formSource.Set("pointer_string", "     pointer value      ")
	formSource.Set("enum", "42")

	sources := map[string]url.Values{
		"url":   urlSource,
		"form":  formSource,
		"query": querySource,
	}

	if err := p.Parse(sources, nil); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "1aa75114-6117-4908-b6ea-0d22ecdd4fc0", s.ID)
	assert.Equal(t, 24, s.Number)
	assert.Equal(t, 0, s.IgnoredInt)
	assert.True(t, s.RequiredBool)
	assert.Nil(t, s.PointerBool)
	assert.Equal(t, "pointer value", *s.PointerString)
	assert.Equal(t, 42, s.Default)
	assert.Equal(t, []int{1, 2, 3}, s.DefaultSlice)
}

func subTestInvalidStruct(t *testing.T) {
	t.Parallel()

	type strct struct {
		ID            string  `from:"url" json:"id" params:"uuid,required"`
		Number        int     `from:"query" json:"number"`
		RequiredBool  bool    `from:"form" json:"required_bool" params:"required"`
		PointerBool   *bool   `from:"form" json:"pointer_bool"`
		PointerString *string `from:"form" json:"pointer_string" params:"trim"`
		Default       int     `from:"form" json:"default" default:"42"`
	}

	sources := map[string]url.Values{
		"url":   url.Values{},
		"query": url.Values{},
		"form":  url.Values{},
	}

	p := params.New(&strct{})
	err := p.Parse(sources, nil)
	assert.Error(t, err)
}

func subTestFieldWithNoSource(t *testing.T) {
	t.Parallel()

	type strct struct {
		ID string `json:"id" params:"uuid"`
	}

	p := params.New(&strct{})
	err := p.Parse(map[string]url.Values{}, nil)
	assert.Error(t, err, "Parse() should have fail")
	assert.Equal(t, "no source set for field ID", err.Error(), "Parse() failed with an unexpected error")
}

func subTestFieldWithUnexistingSource(t *testing.T) {
	t.Parallel()

	type strct struct {
		ID string `from:"somewhere" json:"id" params:"uuid"`
	}

	p := params.New(&strct{})
	err := p.Parse(map[string]url.Values{}, nil)
	assert.Error(t, err, "Parse() should have fail")
	assert.Equal(t, "source somewhere for field ID does not exist", err.Error(), "Parse() failed with an unexpected error")
}

func subTestFieldNotExported(t *testing.T) {
	t.Parallel()

	type strct struct {
		// will fail cause not exported
		id string `from:"url" json:"id" params:"uuid"`
	}

	p := params.New(&strct{})
	err := p.Parse(map[string]url.Values{"url": url.Values{}}, nil)
	assert.Error(t, err, "Parse() should have fail")
	assert.Equal(t, "field id could not be set", err.Error(), "Parse() failed with an unexpected error")
}

func subTestEmbeddedStruct(t *testing.T) {
	t.Parallel()

	type Paginator struct {
		Page    *int `from:"query" json:"page" default:"1"`
		PerPage *int `from:"query" json:"per_page"`
	}

	type strct struct {
		Paginator

		ID string `from:"url" json:"id" params:"uuid,required"`
	}

	s := &strct{}
	p := params.New(s)

	urlSource := url.Values{}
	urlSource.Set("id", "1aa75114-6117-4908-b6ea-0d22ecdd4fc0")

	querySource := url.Values{}
	querySource.Set("page", "24")

	sources := map[string]url.Values{
		"url":   urlSource,
		"query": querySource,
	}

	if err := p.Parse(sources, nil); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "1aa75114-6117-4908-b6ea-0d22ecdd4fc0", s.ID)
	assert.Equal(t, 24, *s.Page)
	assert.Nil(t, s.PerPage)
}

func subTestEmbeddedStructWithCustomValidation(t *testing.T) {
	t.Parallel()

	// sugar
	shouldFail := true

	type strct struct {
		StructWithValidator
	}

	testCases := []struct {
		description string
		params      url.Values
		shouldFail  bool
	}{
		{
			"Trigger a failure",
			url.Values{
				"string":       []string{"value"},
				"true_to_fail": []string{"true"},
			},
			shouldFail,
		},
		{
			"Valid params should work",
			url.Values{
				"string": []string{"value"},
			},
			!shouldFail,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			s := &strct{}
			p := params.New(s)
			sources := map[string]url.Values{
				"query": tc.params,
			}

			err := p.Parse(sources, nil)
			if tc.shouldFail {
				assert.Error(t, err, "Parse() should have failed")
			} else {
				assert.NoError(t, err, "Parse() should have succeed")
			}
		})
	}
}

func subTestCustomValidation(t *testing.T) {
	t.Parallel()

	// sugar
	shouldFail := true

	testCases := []struct {
		description string
		params      url.Values
		shouldFail  bool
	}{
		{
			"Trigger a failure",
			url.Values{
				"string":       []string{"value"},
				"true_to_fail": []string{"true"},
			},
			shouldFail,
		},
		{
			"Valid params should work",
			url.Values{
				"string": []string{"value"},
			},
			!shouldFail,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			s := &StructWithValidator{}
			p := params.New(s)
			sources := map[string]url.Values{
				"query": tc.params,
			}

			err := p.Parse(sources, nil)
			if tc.shouldFail {
				assert.Error(t, err, "Parse() should have failed")
			} else {
				assert.NoError(t, err, "Parse() should have succeed")
			}
		})
	}
}

func subTestExtraction(t *testing.T) {
	t.Parallel()

	cwd, _ := os.Getwd()

	s := struct {
		StructWithValidator
		StringValue   string             `from:"url" json:"string_value"`
		Number        int                `from:"query" json:"number"`
		Bool          bool               `from:"form" json:"bool"`
		PointerBool   *bool              `from:"form" json:"pointer_bool"`
		PointerString *string            `from:"form" json:"pointer_string"`
		PointerNumber *int               `from:"form" json:"pointer_number"`
		Nil           *int               `from:"form" json:"nil"`
		File          *formfile.FormFile `from:"file" json:"file"`
		Stringer      *date.Date         `from:"form" json:"stringer"`
		Ignored       int                `from:"form" json:"-"`
		NoName        string             `from:"form"`
		OmitEmpty     int                `from:"form" json:"omit,omitempty"`
		Slice         []int              `from:"form" json:"slice"`
		NilSlice      []int              `from:"form" json:"nil_slice"`
		EmptySlice    []int              `from:"form" json:"empty_slice"`
		Unknown       string             `json:"unknown"`
	}{
		StringValue:         "String value",
		Number:              42,
		Bool:                true,
		PointerBool:         ptrs.NewBool(false),
		PointerString:       ptrs.NewString("string pointer"),
		PointerNumber:       ptrs.NewInt(24),
		Nil:                 nil,
		File:                testformfile.NewFormFile(t, cwd, "black_pixel.png"),
		Stringer:            date.Today(),
		Ignored:             24,
		NoName:              "not named",
		OmitEmpty:           0,
		Unknown:             "unknown from",
		Slice:               []int{1, 2, 3},
		EmptySlice:          []int{},
		NilSlice:            nil,
		StructWithValidator: StructWithValidator{String: "embeded"},
	}

	p := params.New(&s)
	sources, files := p.Extract()

	// Check file data
	fileData, found := files["file"]
	require.True(t, found, "file should be present")
	assert.NotNil(t, fileData, "fileData should not be nil")
	assert.NotNil(t, fileData.File, "fileData.File should not be nil")
	assert.NotNil(t, fileData.Header, "fileData.header should not be nil")
	assert.Equal(t, "image/png", fileData.Mime)

	// Check url data
	urlValue, found := sources["url"]
	require.True(t, found, "url data should be present")
	assert.Equal(t, s.StringValue, urlValue.Get("string_value"))

	// Check query data
	queryValue, found := sources["query"]
	require.True(t, found, "query data should be present")
	assert.Equal(t, strconv.Itoa(s.Number), queryValue.Get("number"))
	assert.Equal(t, s.StructWithValidator.String, queryValue.Get("string"))

	// Check form data
	formValue, found := sources["form"]
	require.True(t, found, "form data should be present")
	assert.Equal(t, strconv.FormatBool(s.Bool), formValue.Get("bool"))
	assert.Equal(t, strconv.FormatBool(*s.PointerBool), formValue.Get("pointer_bool"))
	assert.Equal(t, *s.PointerString, formValue.Get("pointer_string"))
	assert.Equal(t, strconv.Itoa(*s.PointerNumber), formValue.Get("pointer_number"))
	assert.Equal(t, s.NoName, formValue.Get("NoName"))
	_, exists := formValue["omit"]
	assert.False(t, exists, "omitempty should not exists")
	assert.Empty(t, formValue.Get("nil"))
	assert.Empty(t, formValue.Get("Ignored"))
	d, err := date.New(formValue.Get("stringer"))
	assert.NoError(t, err, "db.NewDate() should have succeed")
	assert.True(t, s.Stringer.Equal(d), "The date changed from %s to %s", s.Stringer, d)
	assert.Equal(t, []string{"1", "2", "3"}, formValue["slice"])
	_, exists = formValue["nil_slice"]
	assert.False(t, exists, "nil_slice should not be part of the output")
	_, exists = formValue["empty_slice"]
	assert.True(t, exists, "empty_slice should be part of the output")

	// Check unknown data
	unknownValue, found := sources["unknown"]
	require.True(t, found, "unknown data should be present")
	assert.Equal(t, s.Unknown, unknownValue.Get("unknown"))
}

func subTestExtractNil(t *testing.T) {
	t.Parallel()

	p := params.New(nil)
	data, files := p.Extract()
	assert.Equal(t, 0, len(files))

	for _, d := range data {
		assert.Equal(t, 0, len(d))
	}
}

func subTestFileUpload(t *testing.T) {
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

			s := strct{}
			p := params.New(&s)
			err := p.Parse(nil, fileHolder)

			if tc.sendNil {
				assert.Error(t, err, "Expected Parse to return an error")
			} else {
				assert.NoError(t, err, "Expected Parse not to return an error")
			}
		})
	}
}
