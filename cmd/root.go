package cmd

import (
	"fmt"
	"os"

	t "github.com/mmmcclimon/toggl-go/internal/toggl"
	"github.com/spf13/cobra"
)

var toggl *t.Toggl

var rootCmd = &cobra.Command{
	Use: "toggl",

	// root command is only interesting as a container
	Run: func(cmd *cobra.Command, args []string) { cmd.Help() },

	// read config, etc.
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		if !cmd.HasParent() {
			// we're the root, only exist for help
			return nil
		}

		toggl = t.NewToggl()
		return toggl.ReadConfig()
	},

	// do not generate a "completion" command, jeez
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

func Execute() {
	// hide all the root help, it's just in the way
	rootCmd.InitDefaultHelpFlag()
	rootCmd.Flags().MarkHidden("help")
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
