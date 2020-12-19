package cmd

import (
	"github.com/jagoe/haste-client-go/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [haste key]",
	Short: "Get a haste from the server",
	Long:  `Get a haste from the server`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config := client.HasteConfig{}
		viper.Unmarshal(&config)

		client.Get(args[0], &config)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// TODO: add --out/-o to write to file
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
