package cmd

import (
	"errors"
	"fmt"
	"os"
	"regexp"
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
	id      string
}

func init() {
	startCmd.Flags().StringVarP(&opts.project, "project", "p", "", "project shortcut for this task")
	startCmd.Flags().StringVarP(&opts.id, "id", "i", "", "jira id for this task")
	rootCmd.AddCommand(startCmd)
}

var likelyId = regexp.MustCompile(`(?i)^[a-z]{3,}-[0-9]+$`)

func runStart(cmd *cobra.Command, args []string) error {
	desc := strings.Join(args, " ")
	if len(desc) == 0 {
		return errors.New("need a description")
	}

	if opts.id != "" || likelyId.MatchString(desc) {
		id := opts.id
		if id == "" {
			id = desc
		}

		startJiraTask(id)
		return nil
	}

	projectId := 0
	if len(opts.project) > 0 {
		projectId = toggl.Config.ProjectShortcuts[opts.project]
	}

	// is this a shortcut
	if strings.HasPrefix(desc, "@") {
		shortcut, ok := toggl.Config.TaskShortcuts[strings.TrimPrefix(desc, "@")]

		if !ok {
			fmt.Printf("could not resolve shortcut %s\n", desc)
			os.Exit(1)
		}

		// no error handling here, just don't mess up your config file, ok
		desc = shortcut["desc"]
		if proj, ok := shortcut["project"]; ok {
			projectId = toggl.Config.ProjectShortcuts[proj]
		}
	}

	// tags
	tag := ""
	words := strings.Split(desc, " ")
	last := words[len(words)-1]

	if strings.HasPrefix(last, "#") {
		tag = strings.TrimPrefix(last, "#")
		desc = strings.Join(words[0:len(words)-1], " ")
	}

	startTask(desc, projectId, tag)
	return nil
}

func startJiraTask(taskId string) {
	c := toggl.Config.NewJiraClient()
	issue := c.GetIssue(taskId)
	startTask(issue.PrettyDescription(), issue.TogglProjectId, "")
}

func startTask(desc string, projectId int, tag string) {
	timer, err := toggl.StartTimer(desc, projectId, tag)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("started timer: %s\n", timer.OnelineDesc())
}
