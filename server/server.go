package server

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jagoe/haste-client-go/config"
)

// Get reads a haste from the provided server
func Get(key string, config *config.GetConfig) (string, error) {
	client := &http.Client{}

	if config.CanProvideClientCertificate() {
		cert, err := tls.LoadX509KeyPair(config.ClientCertificatePath, config.ClientCertificateKeyPath)
		if err != nil {
			return "", fmt.Errorf("Error reading client certificate: %e", err)
		}

		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
		}
	}

	response, err := client.Get(fmt.Sprintf("%s/raw/%s", config.Server, key))
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
