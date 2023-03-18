//go:build jira

package cmd

const JIRA_ENABLED = true

// used by cmd/start
func startJiraTask(taskId string) {
	c := toggl.Config.NewJiraClient()
	issue := c.GetIssue(taskId)
	startTask(issue.PrettyDescription(), issue.TogglProjectId, "")
}
