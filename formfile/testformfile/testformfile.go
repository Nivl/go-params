package testformfile

import (
	"mime/multipart"
	"os"
	"path"
	"testing"

	"github.com/Nivl/go-params/formfile"
	"github.com/Nivl/go-types/filetype"
)

// NewMultipartData is a helper to generate multipart data that can be returned
// by FileHolder.FormFile()
func NewMultipartData(t *testing.T, cwd string, filename string) (*multipart.FileHeader, *os.File) {
	// We create the Head and the file
	header := &multipart.FileHeader{
		Filename: filename,
	}
	filePath := path.Join(cwd, "fixtures", filename)
	file, err := os.Open(filePath)
	if err != nil {
		t.Fatal(err)
	}

	return header, file
}

// NewFormFile is a helper to create a formfile that can be used in a param struct
func NewFormFile(t *testing.T, cwd string, filename string) *formfile.FormFile {
	header, f := NewMultipartData(t, cwd, filename)

	mime, err := filetype.MimeType(f)
	if err != nil {
		t.Fatal(err)
	}

	return &formfile.FormFile{
		File:   f,
		Header: header,
		Mime:   mime,
	}
}
