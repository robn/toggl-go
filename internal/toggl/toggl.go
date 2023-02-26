package toggl

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Toggl struct {
	client http.Client
	cfg    Config
}

type Config struct {
	ApiToken string `toml:"api_token"`
}

func NewToggl() *Toggl {
	return &Toggl{
		cfg:    Config{},
		client: http.Client{},
	}
}

var ErrNoTimer = errors.New("no running timer")

// ReadConfig reads the toggl config file, and returns an error if it can't
// figure out what to read, or if it's not toml
func (t *Toggl) ReadConfig() error {
	filename := os.Getenv("TOGGL_CONFIG_FILE")

	if filename == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("could not determine homedir: %w", err)
		}

		filename = filepath.Join(home, ".togglrc")
	}

	tomlData, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("could not read config file: %w", err)
	}

	_, err = toml.Decode(string(tomlData), &t.cfg)
	if err != nil {
		return fmt.Errorf("bad config file: %w", err)
	}

	return nil
}

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
	return timerFromResponseBody(res.Body)
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

	return timerFromResponseBody(res.Body)
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
