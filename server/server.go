package server

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

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
}

// Get reads a haste from the provided server
func (server HasteServer) Get(key string, client *http.Client) (string, error) {
	if server.ClientCertificatePath != "" && server.ClientCertificateKeyPath != "" {
		if err := server.prepareClientForTLS(client); err != nil {
			return "", err
		}
	}

	response, err := client.Get(fmt.Sprintf("%s/raw/%s", server.URL, key))
	if err != nil {
		return "", fmt.Errorf("Error retrieving haste: %e", err)
	}

	if response.Body != nil {
		defer response.Body.Close()
	}

	if response.StatusCode == 404 {
		return "", fmt.Errorf("No document found: %s", key)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("Error retrieving haste: %e", err)
	}

	return string(body), nil
}

type createHasteResponse struct {
	Key string `json:"key"`
}

// Create a haste on the server
func (server HasteServer) Create(content io.Reader, client *http.Client) (string, error) {
	if server.ClientCertificatePath != "" && server.ClientCertificateKeyPath != "" {
		if err := server.prepareClientForTLS(client); err != nil {
			return "", err
		}
	}

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

func (server HasteServer) prepareClientForTLS(client *http.Client) error {
	cert, err := tls.LoadX509KeyPair(server.ClientCertificatePath, server.ClientCertificateKeyPath)
	if err != nil {
		return fmt.Errorf("Error reading client certificate: %e", err)
	}

	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}

	return nil
}
