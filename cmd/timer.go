package cmd

import (
	"fmt"
	"os"

	t "github.com/mmmcclimon/toggl-go/internal/toggl"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(timerCmd)
}

var timerCmd = &cobra.Command{
	Use:   "timer",
	Short: "what are you doing right now?",
	Run:   runTimer,
}

func runTimer(cmd *cobra.Command, args []string) {
	timer, err := toggl.CurrentTimer()

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

	fmt.Printf("%s so far: %s\n", timer.Duration(), timer.OnelineDesc())
}
