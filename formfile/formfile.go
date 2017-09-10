package formfile

import (
	"mime/multipart"
)

// FormFile represents a file sent
type FormFile struct {
	File   multipart.File
	Header *multipart.FileHeader
	Mime   string
}
