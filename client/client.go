package client

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"github.com/jagoe/haste-client-go/config"
	"github.com/jagoe/haste-client-go/server"
)

// Get retrieves a haste from the server and prints it to STDOUT or into a file
func Get(key string, config *config.GetConfig) {
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
func Create(content io.Reader, config *config.CreateConfig) {
	key, err := server.Create(content, config)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s/%s", config.HasteConfig.Server, key)
}
