package client

import (
	"fmt"
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
