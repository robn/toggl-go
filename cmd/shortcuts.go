package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

var shortcutsCmd = &cobra.Command{
	Use:   "shortcuts",
	Short: "list the things you can start easily",
	Run:   runShortcuts,
}

func init() {
	rootCmd.AddCommand(shortcutsCmd)
}

func runShortcuts(cmd *cobra.Command, args []string) {
	shortcuts := toggl.Config.TaskShortcuts

	titles := maps.Keys(shortcuts)
	sort.Strings(titles)

	length := 0
	for _, title := range titles {
		if len(title) > length {
			length = len(title)
		}
	}

	for _, title := range titles {
		shortcut := shortcuts[title]
		desc := shortcut["desc"]
		project := shortcut["project"]

		// descriptionless task, just use the project as the description
		if desc == "" {
			desc = project
			project = "*taskless*"
		}

		fmt.Printf("@%-*s %s (%s)\n", length+2, title, desc, project)
	}
}
