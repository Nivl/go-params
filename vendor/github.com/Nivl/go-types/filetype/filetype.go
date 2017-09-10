package filetype

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
)

// FileValidator represents a function that can validate a file type
type FileValidator func(r io.ReadSeeker) (bool, error)

// MimeType returns the mimetype of a file
func MimeType(r io.ReadSeeker) (string, error) {
	// DetectContentType needs the first 512 bytes
	bytesNeeded := 512
	buff := make([]byte, bytesNeeded)
	n, err := r.Read(buff)
	if err != nil {
		return "", err
	}

	// we seek back to where we were like we didn't do anything
	if _, err := r.Seek(int64(-n), io.SeekCurrent); err != nil {
		return "", err
	}

	return http.DetectContentType(buff), nil
}

// SHA256Sum returns the SHA256 sum of a reader
func SHA256Sum(r io.ReadSeeker) (string, error) {
	initialPos, err := r.Seek(0, io.SeekCurrent)
	if err != nil {
		return "", err
	}
	// revert the pointer back to its original position
	defer r.Seek(initialPos, io.SeekStart)

	// copy the file to read it's content
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}

	// cast the hash and return it
	hash := fmt.Sprintf("%x", h.Sum(nil))
	return hash, nil
}
