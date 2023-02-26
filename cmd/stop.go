package cmd

import (
	"fmt"
	"os"

	t "github.com/mmmcclimon/toggl-go/internal/toggl"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(stopCmd)
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop doing the thing you're doing",
	Run:   runStop,
}

func runStop(cmd *cobra.Command, args []string) {
	timer, err := toggl.StopCurrentTimer()

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

	fmt.Printf("spent %s: %s\n", timer.Duration(), timer.OnelineDesc())
}
