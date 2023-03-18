package cmd

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:    "config",
	Short:  "dump config and exit (for debugging)",
	Run:    runConfig,
	Hidden: true,
}

func runConfig(cmd *cobra.Command, args []string) {
	cfg := toggl.Config
	spew.Dump(cfg)
}
