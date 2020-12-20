package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/jagoe/haste-client-go/config"
)

// #region Setup
var defaultConfig config.Config = config.Config{
	Server: "hastebin.local",
}

type FakeGetter struct {
	err   error
	haste string
}

func (fake FakeGetter) Get(_ string, _ *http.Client) (string, error) {
	return fake.haste, fake.err
}

type FakeCreator struct {
	err      error
	hasteKey string
}

func (fake FakeCreator) Create(_ io.Reader, _ *http.Client) (string, error) {
	return fake.hasteKey, fake.err
}

// #endregion

func TestGet(t *testing.T) {
	t.Run("should log error", func(t *testing.T) {
		expectedError := "Expected error"
		err := Get("anykey", FakeGetter{err: fmt.Errorf(expectedError)}, bytes.NewBufferString(""))

		if err == nil {
			t.Error("Expected Get to return an error")
		}

		if err.Error() != expectedError {
			t.Errorf("Expected Get to return '%s' as error message, got '%s'", expectedError, err.Error())
		}
	})

	t.Run("should print haste", func(t *testing.T) {
		buffer := bytes.NewBufferString("")
		expectedHaste := "Test haste"

		err := Get("anykey", FakeGetter{haste: expectedHaste}, buffer)

		if err != nil {
			t.Errorf("Expected Get not to return an error, got %e", err)
		}

		haste := buffer.String()
		if haste != expectedHaste {
			t.Errorf("Expected Get to return '%s' as haste, got '%s'", expectedHaste, haste)
		}
	})
}

func TestCreate(t *testing.T) {
	t.Run("should log error", func(t *testing.T) {
		expectedError := "Expected error"
		err := Create(bytes.NewBufferString(""), FakeCreator{err: fmt.Errorf(expectedError)}, "", bytes.NewBufferString(""))

		if err == nil {
			t.Error("Expected Create to return an error")
		}

		if err.Error() != expectedError {
			t.Errorf("Expected Create to return '%s' as error message, got '%s'", expectedError, err.Error())
		}
	})

	t.Run("should print haste URL", func(t *testing.T) {
		buffer := bytes.NewBufferString("")
		serverURL := "hastebin.local"
		hasteKey := "abcdef"
		expectedHasteURL := fmt.Sprintf("%s/%s", serverURL, hasteKey)

		err := Create(bytes.NewBufferString(""), FakeCreator{hasteKey: hasteKey}, serverURL, buffer)

		if err != nil {
			t.Errorf("Expected Create not to return an error, got %e", err)
		}

		hasteURL := buffer.String()
		if hasteURL != expectedHasteURL {
			t.Errorf("Expected Create to return '%s' as haste URL, got '%s'", expectedHasteURL, hasteURL)
		}
	})
}
