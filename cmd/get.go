package cmd

import (
	"github.com/jagoe/haste-client-go/client"
	"github.com/jagoe/haste-client-go/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [haste key or URL]",
	Short: "Get a haste from the server",
	Long: `Get a haste from the configured server (http://hastebin.com by default) by providing a key or directly from
a hastebin server by providing the complete URL (protocol required!).`,
	Example: `haste get oyivuxonema
haste get http://pastebin.com/oyivuxonema`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config := config.GetConfig{}
		viper.Unmarshal(&config)

		out := cmd.Flag("out")
		if out != nil {
			config.OutputPath = out.Value.String()
		}

		client.Get(args[0], &config)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().StringP("out", "o", "", "File path to save the haste")
}
