package server

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/jagoe/haste-client-go/config"
)

// Get reads a haste from the provided server
func Get(key string, config *config.GetConfig) (string, error) {
	client, err := prepareClientForTLS(&config.HasteConfig)
	if err != nil {
		return "", err
	}

	response, err := client.Get(fmt.Sprintf("%s/raw/%s", config.HasteConfig.Server, key))
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
func Create(content io.Reader, config *config.CreateConfig) (string, error) {
	client, err := prepareClientForTLS(&config.HasteConfig)
	if err != nil {
		return "", err
	}

	response, err := client.Post(fmt.Sprintf("%s/documents", config.HasteConfig.Server), "text/plain", content)
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

func prepareClientForTLS(config *config.HasteConfig) (*http.Client, error) {
	client := &http.Client{}

	if !config.CanProvideClientCertificate() {
		return client, nil
	}

	cert, err := tls.LoadX509KeyPair(config.ClientCertificatePath, config.ClientCertificateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("Error reading client certificate: %e", err)
	}

	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}

	return client, nil
}
