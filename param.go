package params

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

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
	// We make sure we are using the right structure to store the file
	if p.Info.Type.String() != "*formfile.FormFile" {
		return fmt.Errorf("the only accepted type for a file is *formfile.FormFile, got %s", p.Info.Type)
	}

	// We parse the tag to get the options
	opts, err := NewOptions(p.Tags)
	if err != nil {
		return err
	}

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

	ff.Mime, err = opts.ValidateFileContent(ff.File)
	if err != nil {
		if err == io.EOF {
			if header.Size == 0 {
				return perror.New(opts.Name, ErrMsgEmptyFile)
			}
			return perror.New(opts.Name, ErrMsgCorruptedFile)
		}
		return err
	}

	p.Value.Set(reflect.ValueOf(ff))
	return nil
}

// SetValue sets the value of the param using the provided source
func (p *Param) SetValue(source url.Values) error {
	// We parse the tag to get the options
	opts, err := NewOptions(p.Tags)
	if err != nil {
		return err
	}
	defaultValue := p.Tags.Get("default")

	// The tag needs to be ignored
	if opts.Ignore {
		return nil
	}

	if opts.Name == "" {
		opts.Name = p.Info.Name
	}

	// if we have a slice we need to treat it differently
	if reflect.Indirect(*p.Value).Kind() == reflect.Slice {
		return p.setSliceValue(source, opts, defaultValue)
	}

	value := opts.ApplyTransformations(source.Get(opts.Name))
	if value == "" {
		value = defaultValue
	}

	_, valueProvided := source[opts.Name]
	sugarIsArrayItem := true
	if err := opts.Validate(value, valueProvided, !sugarIsArrayItem); err != nil {
		return err
	}

	// We now set the value in the struct
	if valueProvided || value != "" {
		// we malloc a zero value if we need to store a pointer
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

// setSliceValue sets the values of the slice param using the provided source
func (p *Param) setSliceValue(source url.Values, opts *Options, defaultValue string) error {
	originalValues, valueProvided := source[opts.Name]
	values := []string{}

	// we make a copy of the original array to keep the original data untouched
	if valueProvided {
		values = make([]string, len(originalValues))
		copy(values, originalValues)
	}

	for i, v := range values {
		values[i] = opts.ApplyTransformations(v)
	}

	// Apply the default value if needed
	if len(values) == 0 && defaultValue != "" {
		values = strings.Split(defaultValue, ",")
	}

	if err := opts.ValidateSlice(values, valueProvided); err != nil {
		return err
	}

	// We now set the values in the struct
	if valueProvided || len(values) > 0 {
		sliceType := p.Value.Type()         // []*date.Date
		sliceStructType := sliceType.Elem() // *date.Date OR date.Date
		isPointer := false

		if sliceStructType.Kind() == reflect.Ptr {
			sliceStructType = sliceStructType.Elem()
			isPointer = true
		}

		// for each type we need to loop over the array of values, cast them
		// to the right type, and assign them to the param
		switch sliceStructType.Kind() {
		case reflect.String:
			var finalValue reflect.Value
			if isPointer {
				finalValue = reflect.MakeSlice(sliceType, len(values), cap(values))
				for i, value := range values {
					v := value
					finalValue.Index(i).Set(reflect.ValueOf(&v))
				}
			} else {
				finalValue = reflect.ValueOf(values)
			}
			p.Value.Set(finalValue)
		case reflect.Bool:
			finalValues := reflect.MakeSlice(sliceType, len(values), cap(values))
			for i, value := range values {
				boolVal, err := strconv.ParseBool(value)
				if err != nil {
					return perror.New(opts.Name, ErrMsgInvalidBoolean)
				}
				var v reflect.Value
				if isPointer {
					v = reflect.ValueOf(&boolVal)
				} else {
					v = reflect.ValueOf(boolVal)
				}
				finalValues.Index(i).Set(v)
			}
			p.Value.Set(finalValues)

		case reflect.Int:
			finalValues := reflect.MakeSlice(sliceType, len(values), cap(values))
			for i, value := range values {
				int64Val, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return perror.New(opts.Name, ErrMsgInvalidInteger)
				}
				intVal := int(int64Val)
				var v reflect.Value
				if isPointer {
					v = reflect.ValueOf(&intVal)
				} else {
					v = reflect.ValueOf(intVal)
				}
				finalValues.Index(i).Set(v)
			}
			p.Value.Set(finalValues)
		case reflect.Struct:
			if _, ok := reflect.New(sliceStructType).Interface().(Scanner); ok {
				slice := reflect.MakeSlice(sliceType, len(values), cap(values))

				for i, value := range values {
					// We need to create a pointer to be able to cast to Scanner
					strct := reflect.New(sliceStructType)
					if scanner, ok := strct.Interface().(Scanner); ok {
						if err := scanner.ScanString(value); err != nil {
							return perror.New(opts.Name, err.Error())
						}
						if !isPointer {
							// Because we want an array of struct, we use Indirect to deference
							// the pointer we have
							strct = reflect.Indirect(strct)
						}
						slice.Index(i).Set(strct)
					}
				}
				p.Value.Set(slice)
			}
		}
	}
	return nil
}
