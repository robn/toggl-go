package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "list the buckets things can go in",
	Run:   runProjects,
}

func init() {
	rootCmd.AddCommand(projectsCmd)
}

func runProjects(cmd *cobra.Command, args []string) {
	projects := toggl.Config.ProjectShortcuts

	shortcuts := maps.Keys(projects)
	sort.Strings(shortcuts)

	for _, sc := range shortcuts {
		fmt.Printf("- %s (%d)\n", sc, projects[sc])
	}
}
