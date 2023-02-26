package toggl

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

type Toggl struct {
	client http.Client
	cfg    Config
}

func NewToggl() *Toggl {
	return &Toggl{
		cfg:    Config{},
		client: http.Client{},
	}
}

var (
	ErrNoTimer = errors.New("no running timer")
)

func urlFor(endpoint string, args ...any) string {
	// https://api.track.toggl.com/api/v9/workspaces/{workspace_id}/time_entries/{time_entry_id}
	return fmt.Sprintf("https://api.track.toggl.com/api/v9"+endpoint, args...)
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
