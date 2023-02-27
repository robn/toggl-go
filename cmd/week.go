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

// This is so goofy: time.Truncate() acts on absolute (roughly, Unix) time,
// and not on the local time, so if it's Monday at 4pm in Philadelphia,
// truncating to 24*hour will give you a time that's Sunday 7pm, rather than
// Monday at midnight, which is what I actually need.
//
// To get around this, we compute now and its current tz offset, truncate
// today, then add back in the offset to correct for it. This is definitely
// broken during DST, when days aren't 24 hours and the offset might change
// over the course of the day, but for this that's fine. Times are hard.
func startOfToday() time.Time {
	now := time.Now()
	_, offset := now.Zone()

	when := now.Truncate(24 * time.Hour)
	return when.Add(time.Duration(-1*offset) * time.Second)
}
