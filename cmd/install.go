package cmd

import "github.com/spf13/cobra"

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "to install the service",
	Long:  "to install the service and write systemd daemon file",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
