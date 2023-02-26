package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:               "toggl",
	Short:             "whatcha up to?",
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		// fmt.Println("hrm")
	},
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
