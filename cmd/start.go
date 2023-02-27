package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start description",
	Short: "start doing a new thing",
	RunE:  runStart,
}

var opts struct {
	project string
}

func init() {
	startCmd.Flags().StringVarP(&opts.project, "project", "p", "", "project shortcut for this task")
	rootCmd.AddCommand(startCmd)
}

func runStart(cmd *cobra.Command, args []string) error {
	desc := strings.Join(args, " ")
	if len(desc) == 0 {
		return errors.New("need a description")
	}

	projectId := 0
	if len(opts.project) > 0 {
		projectId = toggl.Config.ProjectShortcuts[opts.project]
	}

	timer, err := toggl.StartTimer(desc, projectId)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("started timer: %s\n", timer.OnelineDesc())

	return nil
}
