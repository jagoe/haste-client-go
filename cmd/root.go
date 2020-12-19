package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jagoe/haste-client-go/client"
	"github.com/jagoe/haste-client-go/config"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "haste",
	Short: "A hastebin client, written in Go",
	Long:  `A hastebin client that can create hastes from files and STDIN and read hastes from a configurable server.`,
	Example: `echo Test | haste
cat ./file | haste`,
	Run: func(cmd *cobra.Command, args []string) {
		config := config.CreateConfig{}
		viper.Unmarshal(&config)

		client.Create(os.Stdin, &config)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// TODO: set server, client key
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Config file [$HOME/.haste-client-go.yaml]")
	rootCmd.PersistentFlags().StringP("server", "s", "http://hastebin.com", "Server URL")
	rootCmd.PersistentFlags().String("client-cert", "", "Client certificate path")
	rootCmd.PersistentFlags().String("client-cert-key", "", "Client certificate key path")
	viper.BindPFlag("server", rootCmd.PersistentFlags().Lookup("server"))
	viper.BindPFlag("clientCert", rootCmd.PersistentFlags().Lookup("client-cert"))
	viper.BindPFlag("clientCertKey", rootCmd.PersistentFlags().Lookup("client-cert-key"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// TODO: add -f/--file
	// TODO: add -o/--output
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".haste-client-go" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".haste-client-go")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		os.Stderr.WriteString(fmt.Sprintf("Using config file: %s\n", viper.ConfigFileUsed()))
	}
}
