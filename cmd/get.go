package cmd

import (
	"log"

	"github.com/jagoe/haste-client-go/client"
	"github.com/jagoe/haste-client-go/server"
	"github.com/jagoe/haste-client-go/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [haste key or URL]",
	Short: "Get a haste from the server",
	Long: `Get a haste from the configured server (https://hastebin.com by default) by providing a key or directly from
a hastebin server by providing the complete URL (protocol required!).`,
	Example: `haste get oyivuxonema
haste get http://pastebin.com/oyivuxonema`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		server := server.HasteServer{}
		viper.Unmarshal(&server)

		var filepath string
		if cmd.Flag("out") == nil {
			filepath = ""
		} else {
			filepath = cmd.Flag("out").Value.String()
		}

		output, err := client.SetupGetOutput(filepath, client.OsFileOpener{})
		if err != nil {
			log.Fatal(err)
		}

		serverURL, key := util.ParseURL(args[0])
		if serverURL == "" || key == "" {
			// not a valid URL, must be a key - default to configured server and use the provided key
			key = args[0]
		} else {
			// a valid URL - override the configured server and use the parsed key
			server.URL = serverURL
		}

		client.Get(key, server, output)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().StringP("out", "o", "", "File path to save the haste")
}
