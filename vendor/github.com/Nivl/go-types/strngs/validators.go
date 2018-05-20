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

// IsValidSlug returns True if provided text does not contain white characters,
// punctuation, all letters are lower case and only from ASCII range.
// It could contain `-` and `_` but not at the beginning or end of the text.
//
// Taken from https://github.com/gosimple/slug/blob/e9f42fa127660e552d0ad2b589868d403a9be7c6/slug.go#L144
// under Mozilla Public License
func IsValidSlug(text string) bool {
	if text == "" ||
		text[0] == '-' || text[0] == '_' ||
		text[len(text)-1] == '-' || text[len(text)-1] == '_' {
		return false
	}
	for _, c := range text {
		if (c < 'a' || c > 'z') && c != '-' && c != '_' && (c < '0' || c > '9') {
			return false
		}
	}
	return true
}
