package params_test

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	params "github.com/Nivl/go-params"
	"github.com/Nivl/go-params/perror"
	"github.com/Nivl/go-types/date"
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
