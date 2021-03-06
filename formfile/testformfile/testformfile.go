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
func NewMultipartData(t *testing.T, cwd, filename string) (header *multipart.FileHeader, file *os.File) {
	var err error

	// find and open the file
	filePath := path.Join(cwd, "testdata", filename)
	file, err = os.Open(filePath)
	if err != nil {
		t.Fatal(err)
	}

	// we call stat to get data about the file
	stat, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}

	// build a fake header
	header = &multipart.FileHeader{
		Filename: filename,
		Size:     stat.Size(),
	}
	return header, file
}

// NewFormFile is a helper to create a formfile that can be used in a param struct
func NewFormFile(t *testing.T, cwd, filename string) *formfile.FormFile {
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
