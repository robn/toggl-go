package cmd

import (
	"fmt"
	"os"
	"time"

	t "github.com/mmmcclimon/toggl-go/internal/toggl"
	"github.com/spf13/cobra"
)

var weekCmd = &cobra.Command{
	Use:   "week",
	Short: "how's the week been?",
	Run:   runWeek,
}

func init() {
	rootCmd.AddCommand(weekCmd)
}

func runWeek(cmd *cobra.Command, args []string) {
	end := time.Now()
	start := startOfToday()

	// Back up til we hit a Monday (the week starts on Monday, come at me.)
	for start.Local().Weekday() != time.Monday {
		start = start.Add(-24 * time.Hour)
	}

	entries, err := toggl.TimeEntries(start, end)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(entries) == 0 {
		fmt.Println("Nothing logged this week.")
		return
	}

	t.PrintEntryList(entries)
}
