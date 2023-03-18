//go:build !jira

package cmd

const JIRA_ENABLED = false

func startJiraTask(_ string) {
	panic("ended in startJiraTask with no jira")
}
