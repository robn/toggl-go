package toggl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Toggl struct {
	Config Config
	client http.Client
}

func NewToggl() *Toggl {
	return &Toggl{
		Config: Config{},
		client: http.Client{},
	}
}

const UserAgent = "toggl/go v0"

var (
	ErrNoTimer = errors.New("no running timer")
)

func urlFor(endpoint string, args ...any) string {
	// https://api.track.toggl.com/api/v9/workspaces/{workspace_id}/time_entries/{time_entry_id}
	return fmt.Sprintf("https://api.track.toggl.com/api/v9"+endpoint, args...)
}

type startArgs struct {
	Description string   `json:"description"`
	CreatedWith string   `json:"created_with"`
	Start       string   `json:"start"`    // should maybe be time.Time, but wevs
	Duration    int64    `json:"duration"` // silly
	WorkspaceId int      `json:"workspace_id"`
	ProjectId   int      `json:"project_id,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

func (t *Toggl) StartTimer(description string, projectId int, tag string) (*Timer, error) {
	url := urlFor("/workspaces/%d/time_entries", t.Config.WorkspaceId)

	now := time.Now()

	tags := []string{}
	if tag != "" {
		tags = append(tags, tag)
	}

	args := startArgs{
		Description: description,
		CreatedWith: UserAgent,
		Start:       now.UTC().Format("2006-01-02T15:04:05Z"),
		Duration:    now.Unix() * -1,
		WorkspaceId: t.Config.WorkspaceId,
		ProjectId:   projectId,
		Tags:        tags,
	}

	data, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("bogus json: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}

	res, err := t.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("bad post: %w", err)
	}

	defer res.Body.Close()
	return t.timerFromResponseBody(res.Body)
}

func (t *Toggl) CurrentTimer() (*Timer, error) {
	res, err := t.get(urlFor("/me/time_entries/current"))

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	return t.timerFromResponseBody(res.Body)
}

func (t *Toggl) StopCurrentTimer() (*Timer, error) {
	timer, err := t.CurrentTimer()
	if err != nil {
		return nil, err
	}

	url := urlFor("/workspaces/%d/time_entries/%d/stop", timer.WorkspaceId, timer.Id)

	req, err := http.NewRequest(http.MethodPatch, url, nil)
	if err != nil {
		panic(err) // should not happen
	}

	res, err := t.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("bad patch: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Fprintln(os.Stderr, "bad response from patch:")
		dumpResponseAndExit(res)
	}

	return t.timerFromResponseBody(res.Body)
}

func (t *Toggl) AbortCurrentTimer() (*Timer, error) {
	timer, err := t.CurrentTimer()
	if err != nil {
		return nil, err
	}

	url := urlFor("/workspaces/%d/time_entries/%d", timer.WorkspaceId, timer.Id)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		panic(err) // should not happen
	}

	res, err := t.doRequest(req)

	if err != nil {
		return nil, fmt.Errorf("bad abort: %w", err)
	}

	if res.StatusCode != 200 {
		fmt.Fprintln(os.Stderr, "bad response from delete:")
		dumpResponseAndExit(res)
	}

	return timer, nil
}
