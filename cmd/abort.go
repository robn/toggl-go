package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(abortCmd)
}

var abortCmd = &cobra.Command{
	Use:   "abort",
	Short: "actually, you weren't doing that thing after all",
	Run:   runAbort,
}

func runAbort(cmd *cobra.Command, args []string) {
	fmt.Println("abort")
}
