package main

import "github.com/jagoe/haste-client-go/cmd"

var (
	version string
	commit  string
	builtAt string
	builtBy string
)

func main() {
	cmd.Execute()
}
