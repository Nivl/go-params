package ptrs

import (
	"reflect"
)

// Unwrap takes a pointer and return it's value.
// If the pointer is nil, a zero value is returned
func Unwrap(i interface{}) interface{} {
	source := reflect.ValueOf(i)

	// If the provided interface is not a pointer we just return it
	if source.Kind() != reflect.Ptr {
		return i
	}

	// if the interface is a nil pointer we create a variable that contains
	// a zero value of the type of that interface
	// Example: If the interface is a nil string, the logic is the same as doing
	//          var str string
	//          return str
	if source.IsNil() {
		// We get the type of the value (and not the pointer)
		typ := source.Type().Elem()
		return reflect.Zero(typ).Interface()

	}
	// we return the value of the pointer
	return source.Elem().Interface()
}

// UnwrapString takes a pointer and return it's value.
// If the pointer is nil, an empty string is returned
func UnwrapString(val *string) string {
	if val == nil {
		return ""
	}
	return *val
}

// UnwrapInt takes a pointer and return it's value.
// If the pointer is nil, 0 is returned
func UnwrapInt(val *int) int {
	if val == nil {
		return 0
	}
	return *val
}

// UnwrapBool takes a pointer and return it's value.
// If the pointer is nil, false is returned
func UnwrapBool(val *bool) bool {
	if val == nil {
		return false
	}
	return *val
}
