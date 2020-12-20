package server

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// #region Setup

type FakeKeyPairLoader struct {
	cert tls.Certificate
	err  error
}

func (fake FakeKeyPairLoader) LoadX509KeyPair(certFile string, keyFile string) (tls.Certificate, error) {
	return fake.cert, fake.err
}

type TestSettings struct {
	ResponseCode       int
	ResponseBody       string
	KeyPairLoaderCert  tls.Certificate
	KeyPairLoaderError error
}

func prepareTest(settings TestSettings) (HasteServer, *httptest.Server) {
	if settings.ResponseCode == 0 {
		settings.ResponseCode = 200
	}

	fakeServerEndpoint := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(settings.ResponseCode)

		body := settings.ResponseBody + "|" + r.RequestURI
		w.Write([]byte(body))
	}))

	testServer := HasteServer{
		URL:                      fakeServerEndpoint.URL,
		ClientCertificatePath:    "./fake/path",
		ClientCertificateKeyPath: "./fake/path",
		KeyPairLoader:            FakeKeyPairLoader{err: settings.KeyPairLoaderError, cert: settings.KeyPairLoaderCert},
	}

	return testServer, fakeServerEndpoint
}

// #endregion

func TestGet(t *testing.T) {
	t.Run("should return an error if the transport config cannot be set", func(t *testing.T) {
		configError := "Expected error"
		expectedError := fmt.Sprintf("Error reading client certificate: %s", configError)
		server, endpoint := prepareTest(TestSettings{KeyPairLoaderError: fmt.Errorf(configError)})
		defer endpoint.Close()

		_, err := server.Get("anykey", endpoint.Client())

		if err == nil {
			t.Fatalf("Should have returned an error")
		}

		if err.Error() != expectedError {
			t.Fatalf("Should have returned '%s' as error, got '%s'", expectedError, err.Error())
		}
	})

	t.Run("should return an error if the GET request returns a 404 response", func(t *testing.T) {
		key := "abcdef"
		expectedError := fmt.Sprintf("No document found: %s", key)
		server, endpoint := prepareTest(TestSettings{ResponseCode: 404})
		defer endpoint.Close()

		_, err := server.Get(key, endpoint.Client())

		if err == nil {
			t.Fatalf("Should have returned an error")
		}

		if err.Error() != expectedError {
			t.Fatalf("Should have returned '%s' as error, got '%s'", expectedError, err.Error())
		}
	})

	t.Run("should return the response body as text", func(t *testing.T) {
		key := "abcdef"
		haste := "Cool haste, bro!"
		server, endpoint := prepareTest(TestSettings{ResponseBody: haste})
		defer endpoint.Close()

		response, err := server.Get(key, endpoint.Client())

		if err != nil {
			t.Fatalf("Should not have returned an error: %e", err)
		}

		returnedHaste := strings.Split(response, "|")[0]
		if returnedHaste != haste {
			t.Fatalf("Expected response body to be '%s', but got '%s'", haste, returnedHaste)
		}
	})
}

func TestCreate(t *testing.T) {
	t.Run("should return an error if the transport config cannot be set", func(t *testing.T) {
		configError := "Expected error"
		expectedError := fmt.Sprintf("Error reading client certificate: %s", configError)
		server, endpoint := prepareTest(TestSettings{KeyPairLoaderError: fmt.Errorf(configError)})
		defer endpoint.Close()

		_, err := server.Create(bytes.NewBufferString("content"), endpoint.Client())

		if err == nil {
			t.Fatalf("Should have returned an error")
		}

		if err.Error() != expectedError {
			t.Fatalf("Should have returned '%s' as error, got '%s'", expectedError, err.Error())
		}
	})

	t.Run("should return an error if the response JSON is not the expected format", func(t *testing.T) {
		expectedError := "Error when retrieving the haste key: invalid character 'i' looking for beginning of object key string"
		server, endpoint := prepareTest(TestSettings{ResponseBody: "{invalid: json}"})
		defer endpoint.Close()

		_, err := server.Create(bytes.NewBufferString("content"), endpoint.Client())

		if err == nil {
			t.Fatalf("Should have returned an error")
		}

		if err.Error() != expectedError {
			t.Fatalf("Should have returned '%s' as error, got '%s'", expectedError, err.Error())
		}
	})

	t.Run("should return the key of the generated haste", func(t *testing.T) {
		key := "abcdef"
		server, endpoint := prepareTest(TestSettings{ResponseBody: fmt.Sprintf(`{"key": "%s"}`, key)})
		defer endpoint.Close()

		returnedKey, err := server.Create(bytes.NewBufferString("content"), endpoint.Client())

		if err != nil {
			t.Fatalf("Should not have returned an error: %e", err)
		}

		if returnedKey != key {
			t.Fatalf("Should have returned key '%s', got: %s", key, returnedKey)
		}
	})
}
