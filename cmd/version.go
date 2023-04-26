package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	Version     string
	BuildTime   string
	CommitID    string
	BuildBranch string
	GoVersion   string
	OsArch      string
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"ver"},
	Short:   "version",
	Long:    "version and git info...",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\n", Version)
		fmt.Printf("Build Time: %s\n", BuildTime)
		fmt.Printf("Build Branch: %s\n", BuildBranch)
		fmt.Printf("Git Commit: %s\n", CommitID)
		fmt.Printf("Go Version: %s\n", GoVersion)
		fmt.Printf("Os/Arch: %s\n", OsArch)
	},
}
