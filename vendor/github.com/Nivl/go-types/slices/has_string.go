package slices

// HasString looks for a value in a slice, returns true if the value is present,
// false otherwise
func HasString(toFind string, slice []string) bool {
	for _, item := range slice {
		if item == toFind {
			return true
		}
	}
	return false
}
