package server

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// #region Setup

// HasteGetter describes getting hastes from a haste-server instance
type HasteGetter interface {
	Get(key string, client *http.Client) (string, error)
}

// HasteCreator describes creating hastes on a haste-server instance
type HasteCreator interface {
	Create(content io.Reader, client *http.Client) (string, error)
}

// HasteServer provides functionality to interact with a haste-server instance
type HasteServer struct {
	URL                      string `mapstructure:"server"`
	ClientCertificatePath    string `mapstructure:"clientCert"`
	ClientCertificateKeyPath string `mapstructure:"clientCertKey"`

	// KeyPairLoader is not meant to be set manually; call HasteServer.Initialize() instead
	KeyPairLoader X509KeyPairLoader
}

// Initialize configures the HasteServer instance for use in production and can be skipped for tests
func (server HasteServer) Initialize() {
	server.KeyPairLoader = TLSX509KeyPairLoader{}
}

// #endregion

// Get reads a haste from the provided server
func (server HasteServer) Get(key string, client *http.Client) (string, error) {
	tlsConfig, err := getTLSTransportConfig(server.ClientCertificatePath, server.ClientCertificateKeyPath, server.KeyPairLoader)
	if err != nil {
		return "", err
	}

	client.Transport = tlsConfig

	response, err := client.Get(fmt.Sprintf("%s/raw/%s", server.URL, key))
	if err != nil {
		return "", fmt.Errorf("Error retrieving haste: %s", err.Error())
	}

	if response.Body != nil {
		defer response.Body.Close()
	}

	if response.StatusCode == 404 {
		return "", fmt.Errorf("No document found: %s", key)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("Error retrieving haste: %s", err.Error())
	}

	return string(body), nil
}

type createHasteResponse struct {
	Key string `json:"key"`
}

// Create a haste on the server
func (server HasteServer) Create(content io.Reader, client *http.Client) (string, error) {
	tlsConfig, err := getTLSTransportConfig(server.ClientCertificatePath, server.ClientCertificateKeyPath, server.KeyPairLoader)
	if err != nil {
		return "", err
	}

	client.Transport = tlsConfig

	response, err := client.Post(fmt.Sprintf("%s/documents", server.URL), "text/plain", content)
	if err != nil {
		return "", fmt.Errorf("Error creating haste: %e", err)
	}

	if response.Body != nil {
		defer response.Body.Close()
	}

	var haste createHasteResponse
	decoder := json.NewDecoder(response.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&haste); err != nil {
		return "", fmt.Errorf("Error when creating the haste: %e", err)
	}

	return haste.Key, nil
}

// #region Private

// #region Test types & methods
// GetTLSTransportConfig prepares a TLS transport config with the provided certificate and key
// If certificate or key are not specified, an empty (but usable) configuration will be returned.
func getTLSTransportConfig(certFile string, keyFile string, keyPairLoader X509KeyPairLoader) (*http.Transport, error) {
	if certFile == "" || keyFile == "" {
		return &http.Transport{}, nil
	}

	cert, err := keyPairLoader.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("Error reading client certificate: %s", err.Error())
	}

	return &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}, nil
}

// #endregion

// X509KeyPairLoader wraps the behavior of tls.LoadX509KeyPair
type X509KeyPairLoader interface {
	LoadX509KeyPair(certFile, keyFile string) (tls.Certificate, error)
}

// TLSX509KeyPairLoader wraps tls.LoadX509KeyPair
type TLSX509KeyPairLoader struct{}

// LoadX509KeyPair wraps tls.LoadX509KeyPair
func (TLSX509KeyPairLoader) LoadX509KeyPair(certFile string, keyFile string) (tls.Certificate, error) {
	return tls.LoadX509KeyPair(certFile, keyFile)
}

// #endregion
