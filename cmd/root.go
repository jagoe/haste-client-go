package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/jagoe/haste-client-go/client"
	"github.com/jagoe/haste-client-go/server"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	version string
	cfgFile string
)

// NewRootCommand creates a root command which represents the base command when called without any subcommands
func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "haste [file]",
		Short: "A hastebin client, written in Go",
		Long:  fmt.Sprintf("haste v%s\nA hastebin client that can create hastes from files and STDIN and read hastes from a haste-server instance.", version),
		Args:  cobra.MaximumNArgs(1),
		Example: `echo Test | haste
cat ./file | haste
haste ./file`,
		Run: func(cmd *cobra.Command, args []string) {
			displayVersion := false
			versionFlag := cmd.Flag("version")
			if versionFlag != nil {
				var err error
				displayVersion, err = strconv.ParseBool(versionFlag.Value.String())
				if err != nil {
					displayVersion = false
				}
			}

			if displayVersion {
				fmt.Fprintf(cmd.OutOrStdout(), "v%s", version)
				os.Exit(0)
			}

			server := server.MakeHasteServer()
			viper.Unmarshal(&server)

			var filepath string
			if len(args) > 0 {
				filepath = args[0]
			} else {
				filepath = ""
			}
			input, err := client.SetupCreateInput(filepath, client.OsFileOpener{}, cmd.InOrStdin())
			if err != nil {
				fmt.Fprintln(cmd.ErrOrStderr(), err.Error())
				os.Exit(1)
			}

			err = client.Create(input, server, server.URL, cmd.OutOrStdout())
			if err != nil {
				fmt.Fprintln(cmd.ErrOrStderr(), err.Error())
				os.Exit(1)
			}
		},
	}

	initRootCommand(rootCmd)
	addSubCommands(rootCmd)

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd := NewRootCommand()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(rootCmd.ErrOrStderr(), err.Error())
		os.Exit(1)
	}
}

func initRootCommand(rootCmd *cobra.Command) {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "Config file [$HOME/.haste-client-go.yaml]")
	rootCmd.PersistentFlags().StringP("server", "s", "(global) https://hastebin.com", "Server URL")
	rootCmd.PersistentFlags().String("client-cert", "", "(global) Client certificate path")
	rootCmd.PersistentFlags().String("client-cert-key", "", "(global) Client certificate key path")
	viper.BindPFlag("server", rootCmd.PersistentFlags().Lookup("server"))
	viper.BindPFlag("clientCert", rootCmd.PersistentFlags().Lookup("client-cert"))
	viper.BindPFlag("clientCertKey", rootCmd.PersistentFlags().Lookup("client-cert-key"))

	rootCmd.Flags().BoolP("version", "v", false, "Print the version number")
}

func addSubCommands(rootCmd *cobra.Command) {
	rootCmd.AddCommand(NewGetCommand())
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
	if err := viper.ReadInConfig(); err != nil {
		// no config file, no problem
	}
}
