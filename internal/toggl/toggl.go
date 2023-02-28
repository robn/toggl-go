package toggl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sort"
	"time"

	"golang.org/x/exp/maps"
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

func urlFor(endpoint string, args ...any) *url.URL {
	raw := fmt.Sprintf("https://api.track.toggl.com/api/v9"+endpoint, args...)
	u, err := url.Parse(raw)
	if err != nil {
		panic(err)
	}

	return u
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
		Start:       now.UTC().Format(time.RFC3339),
		Duration:    now.Unix() * -1,
		WorkspaceId: t.Config.WorkspaceId,
		ProjectId:   projectId,
		Tags:        tags,
	}

	data, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("bogus json: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url.String(), bytes.NewBuffer(data))
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

func (t *Toggl) ResumeTimer(timer *Timer) (*Timer, error) {
	tag := ""
	if len(timer.Tags) > 0 {
		tag = timer.Tags[0]
	}

	return t.StartTimer(timer.Description, timer.projectId, tag)
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

	loc := urlFor("/workspaces/%d/time_entries/%d/stop", timer.WorkspaceId, timer.Id)

	req, err := http.NewRequest(http.MethodPatch, loc.String(), nil)
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

	loc := urlFor("/workspaces/%d/time_entries/%d", timer.WorkspaceId, timer.Id)

	req, err := http.NewRequest(http.MethodDelete, loc.String(), nil)
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

func (t *Toggl) TimeEntries(start, end time.Time) ([]*Timer, error) {
	loc := urlFor("/me/time_entries")

	params := url.Values{}
	params.Add("start_date", start.UTC().Format(time.RFC3339))
	params.Add("end_date", end.UTC().Format(time.RFC3339))
	loc.RawQuery = params.Encode()

	res, err := t.get(loc)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	return t.timersFromResponseBody(res.Body)
}

func PrintEntryList(entries []*Timer) {
	// we group by task, so we only report "read email" once even if it shows up
	// 10 times in the list
	grouped := map[string][]*Timer{}
	for _, t := range entries {
		k := fmt.Sprintf("%d!%s", t.projectId, t.Description)
		grouped[k] = append(grouped[k], t)
	}

	keys := maps.Keys(grouped)
	sort.Strings(keys)

	var total time.Duration

	for _, k := range keys {
		entries := grouped[k]

		var taskTotal time.Duration
		for _, e := range entries {
			taskTotal += e.Duration()
		}

		total += taskTotal

		fmt.Printf("%5.2fh  %s\n", taskTotal.Hours(), entries[0].OnelineDesc())
	}

	fmt.Println("------")
	fmt.Printf("%5.2fh  total (%s)\n", total.Hours(), total)
}
