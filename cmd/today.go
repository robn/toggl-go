package cmd

import (
	"fmt"
	"os"
	"time"

	t "github.com/mmmcclimon/toggl-go/internal/toggl"
	"github.com/spf13/cobra"
)

var todayCmd = &cobra.Command{
	Use:   "today",
	Short: "what are the things you've done today?",
	Run:   runToday,
}

func init() {
	rootCmd.AddCommand(todayCmd)
}

func runToday(cmd *cobra.Command, args []string) {
	start := startOfToday()
	end := time.Now()

	entries, err := toggl.TimeEntries(start, end)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(entries) == 0 {
		fmt.Println("Nothing logged today.")
		return
	}

	t.PrintEntryList(entries)
}
