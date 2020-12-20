package client

import (
	"fmt"
	"io"
	"os"
)

// FileOpener is an interface that describes opening a file
type FileOpener interface {
	OpenFile(name string, flag int, perm os.FileMode) (*os.File, error)
}

// OsFileOpener wraps os.OpenFile
type OsFileOpener struct{}

// OpenFile wraps os.OpenFile
func (OsFileOpener) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}

// SetupGetOutput prepares the output stream for client.Get based on a filepath that the user did or did not provide
func SetupGetOutput(filepath string, fileOpener FileOpener) (io.Writer, error) {
	if filepath == "" {
		return os.Stdout, nil
	}

	file, err := fileOpener.OpenFile(filepath, os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("Error creating output file: %s", err.Error())
	}

	return file, nil

}
