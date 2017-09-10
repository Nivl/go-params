package ptrs

// NewBool takes a bool and returns a pointer to a new bool of the same value
func NewBool(v bool) *bool {
	b := v
	return &b
}

// NewInt takes an int and returns a pointer to a new int of the same value
func NewInt(v int) *int {
	i := v
	return &i
}

// NewString takes a string and returns a pointer to a new string of the same value
func NewString(s string) *string {
	return &s
}
