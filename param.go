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
	Value *reflect.Value
	Info  *reflect.StructField
	Tags  *reflect.StructTag
}

// SetFile sets the value of the param using the provided source to find the file
func (p *Param) SetFile(source formfile.FileHolder) error {
	// We parse the tag to get the options
	opts := NewOptions(p.Tags)

	// The tag needs to be ignored
	if opts.Ignore {
		return nil
	}

	if opts.Name == "" {
		opts.Name = p.Info.Name
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
	if p.Info.Type.String() != "*formfile.FormFile" {
		return fmt.Errorf("the only accepted type for a file is *formfile.FormFile, got %s", p.Info.Type)
	}

	ff.Mime, err = opts.ValidateFileContent(ff.File)
	if err != nil {
		return err
	}

	p.Value.Set(reflect.ValueOf(ff))
	return nil
}

// SetValue sets the value of the param using the provided source
func (p *Param) SetValue(source url.Values) error {
	// We parse the tag to get the options
	opts := NewOptions(p.Tags)
	defaultValue := p.Tags.Get("default")

	// The tag needs to be ignored
	if opts.Ignore {
		return nil
	}

	if opts.Name == "" {
		opts.Name = p.Info.Name
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
		if p.Value.Kind() == reflect.Ptr {
			val := reflect.New(p.Value.Type().Elem())
			p.Value.Set(val)
		}

		field := reflect.Indirect(*p.Value)
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
			if scanner, ok := p.Value.Interface().(Scanner); ok {
				if err := scanner.ScanString(value); err != nil {
					return perror.New(opts.Name, err.Error())
				}
			}
		}
	}
	return nil
}
