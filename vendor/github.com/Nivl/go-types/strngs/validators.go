package strngs

import (
	"net/url"
	"regexp"
	"strings"
)

// IsValidUUID checks if a string is a valid UUID V4
func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[8|9|aA|bB][a-f0-9]{3}-[a-f0-9]{12}")
	return r.MatchString(uuid)
}

// IsValidURL checks if a string is a valid URL
func IsValidURL(uri string) bool {
	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return false
	}
	// we only accept http(s)
	return strings.HasPrefix(u.Scheme, "http")
}

// IsValidEmail checks if a string is a valid email
// We only check for anything@anything.anything
func IsValidEmail(email string) bool {
	r := regexp.MustCompile(`.+@.+\..+`)
	return r.MatchString(email)
}
