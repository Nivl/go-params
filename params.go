package params

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/Nivl/go-params/formfile"
	"github.com/Nivl/go-params/perror"
)

// Params is a struct used to parse and extract params from an other struct
type Params struct {
	data interface{}
}

// New creates a new Params object from a struct
func New(data interface{}) *Params {
	return &Params{
		data: data,
	}
}

// Parse fills the paramsStruct using the provided sources
func (p *Params) Parse(sources map[string]url.Values, fileHolder formfile.FileHolder) error {
	paramList := reflect.Indirect(reflect.ValueOf(p.data))
	err := p.parseRecursive(paramList, sources, fileHolder)
	if err != nil {
		return err
	}

	// If there's a custom validator we'll use it
	if validator, ok := p.data.(CustomValidation); ok {
		isValid, field, err := validator.IsValid()
		if !isValid {
			return perror.New(field, err.Error())
		}
	}

	return nil
}

func (p *Params) parseRecursive(paramList reflect.Value, sources map[string]url.Values, fileHolder formfile.FileHolder) error {
	nbParams := paramList.NumField()
	for i := 0; i < nbParams; i++ {
		value := paramList.Field(i)
		info := paramList.Type().Field(i)
		tags := info.Tag

		// We make sure we can update the value of field
		if !value.CanSet() {
			return fmt.Errorf("field %s could not be set", info.Name)
		}

		// Handle embedded struct
		if value.Kind() == reflect.Struct && info.Anonymous {
			p.parseRecursive(value, sources, fileHolder)

			// If there's a custom validator we'll use it
			// Here we vall Addr() to make sure we get a pointer to the
			// struct. If we don't use a pointer and the IsValid() method
			// uses a pointer, the conversion will fail
			if validator, ok := value.Addr().Interface().(CustomValidation); ok {
				isValid, field, err := validator.IsValid()
				if !isValid {
					return perror.New(field, err.Error())
				}
			}

			continue
		}

		// We control the source of the param. If nothing is provided, we take from the URL
		paramLocation := strings.ToLower(tags.Get("from"))
		if paramLocation == "" {
			return fmt.Errorf("no source set for field %s", info.Name)
		}

		param := &Param{
			Value: &value,
			Info:  &info,
			Tags:  &tags,
		}

		// the "file" source is a special case as it's not part of the sources object
		if paramLocation == "file" {
			if err := param.SetFile(fileHolder); err != nil {
				return err
			}
		} else {
			source, found := sources[paramLocation]
			if !found {
				return fmt.Errorf("source %s for field %s does not exist", paramLocation, info.Name)
			}

			if err := param.SetValue(source); err != nil {
				return err
			}
		}
	}

	return nil
}

// Extract extracts the data from the paramsStruct and returns them
// as a map of url.Values
func (p *Params) Extract() (map[string]url.Values, map[string]*formfile.FormFile) {
	sources := map[string]url.Values{}
	files := map[string]*formfile.FormFile{}

	if p.data == nil {
		return sources, files
	}

	paramList := reflect.Indirect(reflect.ValueOf(p.data))
	p.extractRecursive(paramList, sources, files)
	return sources, files
}

func (p *Params) extractRecursive(paramList reflect.Value, sources map[string]url.Values, files map[string]*formfile.FormFile) {
	nbParams := paramList.NumField()
	for i := 0; i < nbParams; i++ {
		value := paramList.Field(i)
		info := paramList.Type().Field(i)
		tags := info.Tag

		// skip the nil pointers
		if value.Kind() == reflect.Ptr && value.IsNil() {
			continue
		}

		// We get the name from the json tag
		fieldName := ""
		jsonOpts := strings.Split(tags.Get("json"), ",")
		if len(jsonOpts) > 0 {
			if jsonOpts[0] == "-" {
				continue
			}
			fieldName = jsonOpts[0]
		}
		if fieldName == "" {
			fieldName = info.Name
		}
		// if the field has the omitempty option we want to honor it
		omitempty := false
		if len(jsonOpts) > 1 {
			for _, opt := range jsonOpts {
				if opt == "omitempty" {
					omitempty = true
				}
			}
		}

		// Handle embedded struct
		if reflect.Indirect(value).Kind() == reflect.Struct && info.Anonymous {
			p.extractRecursive(value, sources, files)
			continue
		}

		// We get the source type (url, query, form, ...) and add the value
		sourceType := strings.ToLower(tags.Get("from"))
		if sourceType == "" {
			sourceType = "unknown"
		}

		if _, found := sources[sourceType]; !found {
			sources[sourceType] = url.Values{}
		}

		// Special cases for files
		if info.Type.String() == "*formfile.FormFile" {
			files[fieldName] = value.Interface().(*formfile.FormFile)
			continue
		}

		field := reflect.Indirect(value)
		valueStr := ""
		isZeroValue := false
		switch field.Kind() {
		case reflect.Slice:
			totalElems := value.Len()
			for i := 0; i < totalElems; i++ {
				stringValue := fmt.Sprintf("%v", value.Index(i).Interface())
				sources[sourceType].Add(fieldName, stringValue)
			}
			// special case so we return right away
			continue
		default:
			// we cast the value to string (works with any stringers)
			valueStr = fmt.Sprintf("%v", field.Interface())
			isZeroValue = reflect.Zero(field.Type()).Interface() == field.Interface()
		}

		// if the omitempty option is set, we wont set any zero value
		if !omitempty || (omitempty && !isZeroValue) {
			sources[sourceType].Set(fieldName, valueStr)
		}
	}
}
