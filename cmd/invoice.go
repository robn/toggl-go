package cmd

import (
	"fmt"
	"os"
	"time"

	t "github.com/mmmcclimon/toggl-go/internal/toggl"
	"github.com/spf13/cobra"
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

	t.PrintEntryList(entries)
}
