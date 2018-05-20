package slices

// HasInt looks for a value in a slice, returns true if the value is present,
// false otherwise
func HasInt(toFind int, slice []int) bool {
	for _, item := range slice {
		if item == toFind {
			return true
		}
	}
	return false
}
