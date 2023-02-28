package cmd

import (
	"fmt"
	"os"
	"time"

	t "github.com/mmmcclimon/toggl-go/internal/toggl"
	"github.com/spf13/cobra"
)

var toggl *t.Toggl

var rootCmd = &cobra.Command{
	Use: "toggl",

	// read config, etc.
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		if !cmd.HasParent() {
			// we're the root, only exist for help
			return nil
		}

		toggl = t.NewToggl()
		return toggl.ReadConfig()
	},

	// do not generate a "completion" command, jeez
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

func Execute() {
	// hide all the root help, it's just in the way
	rootCmd.InitDefaultHelpFlag()
	_ = rootCmd.Flags().MarkHidden("help") // will not fail; we know --help exists
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// This is so goofy: time.Truncate() acts on absolute (roughly, Unix) time,
// and not on the local time, so if it's Monday at 4pm in Philadelphia,
// truncating to 24*hour will give you a time that's Sunday 7pm, rather than
// Monday at midnight, which is what I actually need.
//
// To get around this, we do a stupid hack of reparsing the date-only format
// in a local time zone.
func startOfToday() time.Time {
	now := time.Now() // always in Local zone
	format := time.DateOnly
	midnight, _ := time.ParseInLocation(format, now.Format(format), time.Local)
	return midnight
}
