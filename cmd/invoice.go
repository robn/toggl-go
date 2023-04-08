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
	grouped := map[string][]*t.Timer{}
	for _, t := range entries {
		entryStart := t.Start
		// Back up til we hit a Monday
		for entryStart.Weekday() != time.Monday {
			entryStart = entryStart.Add(-24 * time.Hour)
		}
		k := entryStart.Format("2006-01-02")
		grouped[k] = append(grouped[k], t)
	}

	keys := maps.Keys(grouped)
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Println(k)
		t.PrintEntryList(grouped[k])
	}
}
