package cmd

import (
	"fmt"
	"os"

	t "github.com/mmmcclimon/toggl-go/internal/toggl"
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
	timer, err := toggl.AbortCurrentTimer()

	if err != nil {
		switch err {
		case t.ErrNoTimer:
			fmt.Println("You don't have a running timer!")
			return
		default:
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	fmt.Printf("aborted timer: %s\n", timer.OnelineDesc())
}
