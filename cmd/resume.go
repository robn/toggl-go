package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var resumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "restart the last thing you were doing",
	Run:   runResume,
}

func init() {
	rootCmd.AddCommand(resumeCmd)
}

func runResume(cmd *cobra.Command, args []string) {
	start := time.Now().Add(-6 * time.Hour)
	end := time.Now()

	entries, err := toggl.TimeEntries(start, end)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(entries) == 0 {
		fmt.Println("I dunno what you were last up to, sorry.")
		return
	}

	last := entries[0]

	timer, err := toggl.ResumeTimer(last)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("started timer: %s\n", timer.OnelineDesc())
}
