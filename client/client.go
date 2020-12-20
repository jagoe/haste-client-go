package client

import (
	"fmt"
	"io"
	"net/http"

	"github.com/jagoe/haste-client-go/server"
)

// Get retrieves a haste from the server and prints it to STDOUT or into a file
func Get(key string, getter server.HasteGetter, out io.Writer) error {
	haste, err := getter.Get(key, &http.Client{})
	if err != nil {
		return err
	}

	fmt.Fprint(out, haste)
	return nil
}

// Create a new haste on the server and print an identifier to STDOUT
func Create(input io.Reader, creator server.HasteCreator, serverURL string, out io.Writer) error {
	key, err := creator.Create(input, &http.Client{})
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "%s/%s", serverURL, key)
	return nil
}
