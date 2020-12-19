package client

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"github.com/jagoe/haste-client-go/config"
	"github.com/jagoe/haste-client-go/server"
)

// Get retrieves a haste from the server and prints it to STDOUT or into a file
func Get(keyOrURL string, config *config.GetConfig) {
	serverURL, key := parseHasteURL(keyOrURL)
	if serverURL == "" || key == "" {
		// default to configured server and use the provided key
		key = keyOrURL
	} else {
		// override the configured server
		config.HasteConfig.Server = serverURL
	}

	haste, err := server.Get(key, config)
	if err != nil {
		log.Fatal(err)
	}

	if !config.ShouldSaveAsFile() {
		fmt.Println(haste)
		return
	}

	ioutil.WriteFile(config.OutputPath, []byte(haste), 0770)
}

// Create a new haste on the server and print an identifier to STDOUT
func Create(filepath string, config *config.CreateConfig) {
	var input io.Reader
	if filepath != "" {
		file, err := os.Open(filepath)
		if err != nil {
			log.Fatal(err)
		}

		input = file
	} else {
		input = os.Stdin
	}

	key, err := server.Create(input, config)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s/%s", config.HasteConfig.Server, key)
}

func parseHasteURL(possibleURL string) (string, string) {
	r := regexp.MustCompile(`(.*?//.*?)/(.*?)$`)
	match := r.FindStringSubmatch(possibleURL)

	if len(match) < 3 {
		// no match
		return "", ""
	}

	return match[1], match[2]
}
