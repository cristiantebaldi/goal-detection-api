package cmd

import (
	"fmt"
	"strings"

	"github.com/cristiantebaldi/goal-detection-api/pkg/rest"
	"github.com/cristiantebaldi/goal-detection-api/pkg/validator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.SetDefault("server.http.port", "8080")

		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				fmt.Println("Config file not found, using environment variables...")
			} else {
				panic(fmt.Errorf("fatal error config file: %w", err))
			}
		}

		validator.Initialize(validator.GetDefaultConfig())
		rest.Initialize()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
