package client

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// HasteConfig represents the configuration values provided by a config YAML, ENV or flags
type HasteConfig struct {
	Server                   string `mapstructure:"server"`
	ClientCertificatePath    string `mapstructure:"clientCert"`
	ClientCertificateKeyPath string `mapstructure:"clientCertKey"`
}

// CanProvideClientCertificate checks if the config contains all necessary information to provide a client certificate
func (config *HasteConfig) CanProvideClientCertificate() bool {
	return len(config.ClientCertificatePath) > 0 && len(config.ClientCertificateKeyPath) > 0
}

// Get retrieves a haste from the server and prints it to STDOUT
func Get(key string, config *HasteConfig) {
	client := &http.Client{}

	if config.CanProvideClientCertificate() {
		cert, err := tls.LoadX509KeyPair(config.ClientCertificatePath, config.ClientCertificateKeyPath)
		if err != nil {
			log.Fatalf("Error reading client certificate: %e", err)
		}

		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
		}
	}

	response, err := client.Get(fmt.Sprintf("%s/raw/%s", config.Server, key))
	if err != nil {
		log.Fatalf("Error retrieving haste: %e", err)
	}

	if response.Body != nil {
		defer response.Body.Close()
	}

	if response.StatusCode == 404 {
		os.Stderr.WriteString(fmt.Sprintf("No document found: %s\n", key))
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Error retrieving haste: %e", err)
	}

	haste := string(body)

	fmt.Println(haste)
}
