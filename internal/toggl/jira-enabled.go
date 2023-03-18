//go:build jira

package toggl

import (
	"fmt"
	"os"

	"github.com/mmmcclimon/toggl-go/internal/jira"
)

func (cfg *Config) NewJiraClient() *jira.Client {
	jiraConf := cfg.JiraConfig
	if jiraConf == nil || jiraConf.URL == "" || jiraConf.ConsumerKey == "" {
		fmt.Fprintln(os.Stderr, "Unable to create jira client without jira config")
		os.Exit(1)
	}

	return jira.NewClient(jiraConf)
}
