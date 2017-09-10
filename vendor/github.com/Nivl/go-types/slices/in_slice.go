package slices

import (
	"errors"
	"reflect"
)

// InSlice looks for a value in a slice, returns true if the value is present,
// false otherwise
func InSlice(slice interface{}, val interface{}) (bool, error) {
	switch reflect.TypeOf(slice).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(slice)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) {
				return true, nil
			}
		}
	default:
		return false, errors.New("not a slice")
	}
	return false, nil
}
