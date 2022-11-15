package cmd

import (
	"github.com/spf13/cobra"
)

var (
	configPath string

	rootCmd = &cobra.Command{
		Use: "user",
	}
)

func init() {
	rootCmd.AddCommand(appCmd)
}

// getConfigPath need for providing
func getConfigPath() string {
	return configPath
}

func Execute() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "config file")
	_ = rootCmd.MarkPersistentFlagRequired("config")

	_ = rootCmd.Execute()
}
