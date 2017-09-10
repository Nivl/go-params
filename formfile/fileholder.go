//go:generate mockgen -destination mockformfile/fileholder.go -package mockformfile github.com/Nivl/go-params/formfile FileHolder

package formfile

import "mime/multipart"

// FileHolder represents an interface to fetch files from a key
type FileHolder interface {
	FormFile(key string) (multipart.File, *multipart.FileHeader, error)
}
