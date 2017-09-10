package params

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"

	"github.com/Nivl/go-params/perror"

	"github.com/Nivl/go-params/formfile"
)

// Param represents a struct param
type Param struct {
	value  *reflect.Value
	info   *reflect.StructField
	tags   *reflect.StructTag
	source url.Values
}

// NewParamFromStructValue creates a param using a struct value
func NewParamFromStructValue(paramList *reflect.Value, paramPos int) *Param {
	value := paramList.Field(paramPos)
	info := paramList.Type().Field(paramPos)
	tags := info.Tag

	return &Param{
		value: &value,
		info:  &info,
		tags:  &tags,
	}
}

// SetFile sets the value of the param using the provided source to find the file
func (p *Param) SetFile(source formfile.FileHolder) error {
	// We parse the tag to get the options
	opts := NewOptions(p.tags)

	// The tag needs to be ignored
	if opts.Ignore {
		return nil
	}

	if opts.Name == "" {
		opts.Name = p.info.Name
	}

	file, header, err := source.FormFile(opts.Name)
	if err != nil {
		// if the file is missing it's ok as long as it's not required
		if err == http.ErrMissingFile {
			if opts.Required {
				return perror.New(opts.Name, ErrMsgMissingParameter)
			}
			// if there's no file and it's not required, then we're done
			return nil
		}
		return err
	}

	ff := &formfile.FormFile{
		File:   file,
		Header: header,
	}
	if p.info.Type.String() != "*formfile.FormFile" {
		return fmt.Errorf("the only accepted type for a file is *formfile.FormFile, got %s", p.info.Type)
	}

	ff.Mime, err = opts.ValidateFileContent(ff.File)
	if err != nil {
		return err
	}

	p.value.Set(reflect.ValueOf(ff))
	return nil
}

// SetValue sets the value of the param using the provided source
func (p *Param) SetValue(source url.Values) error {
	// We parse the tag to get the options
	opts := NewOptions(p.tags)
	defaultValue := p.tags.Get("default")

	// The tag needs to be ignored
	if opts.Ignore {
		return nil
	}

	if opts.Name == "" {
		opts.Name = p.info.Name
	}

	value := opts.ApplyTransformations(source.Get(opts.Name))
	if value == "" {
		value = defaultValue
	}

	_, valueProvided := source[opts.Name]
	if err := opts.Validate(value, valueProvided); err != nil {
		return err
	}

	// We now set the value in the struct
	if valueProvided || value != "" {
		if p.value.Kind() == reflect.Ptr {
			val := reflect.New(p.value.Type().Elem())
			p.value.Set(val)
		}

		field := reflect.Indirect(*p.value)
		switch field.Kind() {
		case reflect.Bool:
			v, err := strconv.ParseBool(value)
			if err != nil {
				return perror.New(opts.Name, ErrMsgInvalidBoolean)
			}
			field.SetBool(v)
		case reflect.String:
			field.SetString(value)
		case reflect.Int:
			v, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return perror.New(opts.Name, ErrMsgInvalidInteger)
			}
			field.SetInt(v)
		case reflect.Struct:
			if scanner, ok := p.value.Interface().(Scanner); ok {
				if err := scanner.ScanString(value); err != nil {
					return perror.New(opts.Name, err.Error())
				}
			}
		}
	}
	return nil
}
