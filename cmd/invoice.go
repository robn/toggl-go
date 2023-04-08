package cmd

import (
	"fmt"
	"os"
	"sort"
	"time"

	t "github.com/mmmcclimon/toggl-go/internal/toggl"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

var invoiceCmd = &cobra.Command{
	Use:   "invoice",
	Short: "generate an invoice for last month",
	Run:   runInvoice,
}

func init() {
	rootCmd.AddCommand(invoiceCmd)
}

func runInvoice(cmd *cobra.Command, args []string) {
	now := time.Now();
	year, month, _ := now.Date()
	start := time.Date(year, month-1, 1, 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 1, 0).Add(-time.Second)

	entries, err := toggl.TimeEntries(start, end)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(entries) == 0 {
		fmt.Println("Nothing logged last month.")
		return
	}

	// group by week
	byWeek := map[string][]*t.Timer{}
	for _, t := range entries {
		entryStart := t.Start
		// Back up til we hit a Monday
		for entryStart.Weekday() != time.Monday {
			entryStart = entryStart.Add(-24 * time.Hour)
		}
		k := entryStart.Format("2006-01-02")
		byWeek[k] = append(byWeek[k], t)
	}

	keys := maps.Keys(byWeek)
	sort.Strings(keys)

	var invoiceTotal time.Duration

	for _, weekDate := range keys {
		fmt.Printf("\nweek of: %s\n", weekDate)

		weekEntries := byWeek[weekDate]

		// group by project
		byProject := map[int][]*t.Timer{}
		for _, t := range weekEntries {
			byProject[t.ProjectId] =
			    append(byProject[t.ProjectId], t)
		}

		keys := maps.Keys(byProject)
		sort.Ints(keys)

		var weekTotal time.Duration

		for _, projectId := range keys {
			entries := byProject[projectId]

			var projectTotal time.Duration
			for _, e := range entries {
				projectTotal += e.Duration()
			}
			weekTotal += projectTotal

			fmt.Printf("%10d  %5.2fh  %s\n", projectId,
			    projectTotal.Hours(), entries[0].OnelineDesc())
		}

		fmt.Println("            ------")
		fmt.Printf("            %5.2fh\n", weekTotal.Hours())

		invoiceTotal += weekTotal
	}

	fmt.Printf("\ninvoice total: %5.2fh\n", invoiceTotal.Hours())
}
