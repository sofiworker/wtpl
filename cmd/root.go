package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"wtpl/conf"
	"wtpl/pkg/logger"
	"wtpl/pkg/os"
)

var (
	mode string
)

var rootCmd = &cobra.Command{
	Use: os.MustGetAppName(),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	cobra.OnInitialize(func() {
		conf.InitConfig()
		logger.Init(logger.WithLogEncoding(logger.ConsoleEncoder))
	})

	rootCmd.PersistentFlags().StringVarP(&mode, "mode", "", "prod", "can use prod or dev")
	_ = viper.BindPFlag("app.mode", rootCmd.PersistentFlags().Lookup("mode"))

	rootCmd.AddCommand(versionCmd)

	rootCmd.AddCommand(serverCmd)

	rootCmd.AddCommand(devCmd)

	rootCmd.AddCommand(installCmd)
}

func Run() error {
	return rootCmd.Execute()
}
