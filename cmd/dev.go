package cmd

import "github.com/spf13/cobra"

var devCmd = &cobra.Command{
	Use:    "dev",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
