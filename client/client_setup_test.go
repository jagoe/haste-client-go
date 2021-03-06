package client

import (
	"fmt"
	"os"
	"testing"
)

type FakeFileOpener struct {
	err  error
	file *os.File
}

func (fake FakeFileOpener) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	if fake.err != nil {
		return nil, fake.err
	}

	return fake.file, nil
}

func (fake FakeFileOpener) Open(name string) (*os.File, error) {
	if fake.err != nil {
		return nil, fake.err
	}

	return fake.file, nil
}

func TestSetupGetOutput(t *testing.T) {
	t.Run("should return STDOUT if filepath is empty", func(t *testing.T) {
		out, err := SetupGetOutput("", nil, os.Stdout)

		if err != nil {
			t.Fatalf("Should not have returned error: %s", err.Error())
		}

		if out != os.Stdout {
			t.Fatalf("Output is not STDOUT, but %v", out)
		}
	})

	t.Run("should return error if filepath is invalid", func(t *testing.T) {
		fileError := "Expected error"
		expectedError := fmt.Sprintf("Error creating output file: %s", fileError)
		fakeFileOpener := FakeFileOpener{err: fmt.Errorf(fileError)}
		_, err := SetupGetOutput("invalid", fakeFileOpener, nil)

		if err == nil {
			t.Fatalf("Should have returned error")
		}

		if err.Error() != expectedError {
			t.Fatalf("Expected '%s', but got '%s'", expectedError, err.Error())
		}
	})

	t.Run("should return file if filepath is valid", func(t *testing.T) {
		expectedFile := &os.File{}
		fakeFileOpener := FakeFileOpener{file: expectedFile}
		file, err := SetupGetOutput("valid", fakeFileOpener, nil)

		if err != nil {
			t.Fatalf("Should not have returned error: %s", err.Error())
		}

		if file != expectedFile {
			t.Fatalf("Expected '%v', but got '%v'", expectedFile, file)
		}
	})
}

func TestSetupCreateInput(t *testing.T) {
	t.Run("should return STDIN if filepath is empty", func(t *testing.T) {
		out, err := SetupCreateInput("", nil, os.Stdin)

		if err != nil {
			t.Fatalf("Should not have returned error: %s", err.Error())
		}

		if out != os.Stdin {
			t.Fatalf("Output is not STDIN, but %v", out)
		}
	})

	t.Run("should return error if filepath is invalid", func(t *testing.T) {
		fileError := "Expected error"
		expectedError := fmt.Sprintf("Error reading input file: %s", fileError)
		fakeFileOpener := FakeFileOpener{err: fmt.Errorf(fileError)}
		_, err := SetupCreateInput("invalid", fakeFileOpener, nil)

		if err == nil {
			t.Fatalf("Should have returned error")
		}

		if err.Error() != expectedError {
			t.Fatalf("Expected '%s', but got '%s'", expectedError, err.Error())
		}
	})

	t.Run("should return file if filepath is valid", func(t *testing.T) {
		expectedFile := &os.File{}
		fakeFileOpener := FakeFileOpener{file: expectedFile}
		file, err := SetupCreateInput("valid", fakeFileOpener, nil)

		if err != nil {
			t.Fatalf("Should not have returned error: %s", err.Error())
		}

		if file != expectedFile {
			t.Fatalf("Expected '%v', but got '%v'", expectedFile, file)
		}
	})
}
