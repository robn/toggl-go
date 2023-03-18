package jira

import (
	"context"
	"fmt"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/dghubble/oauth1"
)

// lol, obviously this is _very_ specific to my use case.
const (
	epicLinkFieldName = "customfield_10857"
	bfTaskName        = "Build Failure"
)

type Client struct {
	*jira.Client
	cfg *Config
}

func NewClient(cfg *Config) *Client {
	cfg.ensureKeyLoaded()

	ctx := context.Background()

	oauthCfg := oauth1.Config{
		ConsumerKey: cfg.ConsumerKey,
		Signer: &oauth1.RSASigner{
			PrivateKey: cfg.privKey,
		},
	}

	httpClient := oauthCfg.Client(ctx, oauth1.NewToken(cfg.AccessToken, cfg.AccessSecret))
	jiraClient, err := jira.NewClient(httpClient, cfg.URL)
	maybeDie(err != nil, "unable to create new JIRA client. %v", err)

	return &Client{jiraClient, cfg}
}

type Issue struct {
	Raw            *jira.Issue
	Key            string
	Title          string
	Kind           string
	Epic           string
	TogglProjectId int
}

func (c *Client) GetIssue(id string) Issue {
	issue, _, err := c.Issue.Get(id, nil)
	maybeDie(err != nil, "could not get issue %s: %v", id, err)
	return c.issueFromRaw(issue)
}

func (c *Client) issueFromRaw(issue *jira.Issue) Issue {
	fields := issue.Fields

	epic, _ := fields.Unknowns.String(epicLinkFieldName)

	// This is a Toggl project id, which we'll fill in based on the issue.
	// Again, this is obviously very specific to me.
	var projectId int

	if fromProject, ok := c.cfg.Projects[epic]; ok {
		projectId = fromProject
	} else if fields.Type.Name == bfTaskName {
		projectId = c.cfg.Projects["BF_DEFAULT"]
	} else {
		projectId = c.cfg.Projects["DEFAULT"]
	}

	return Issue{
		Raw:            issue,
		Key:            issue.Key,
		Title:          fields.Summary,
		Kind:           fields.Type.Name,
		Epic:           epic,
		TogglProjectId: projectId,
	}
}

func (iss *Issue) PrettyDescription() string {
	lcKey := strings.ToLower(iss.Key)

	if iss.Kind == bfTaskName {
		return lcKey
	}

	return fmt.Sprintf("%s: %s", lcKey, iss.Title)
}
