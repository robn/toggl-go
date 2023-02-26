package toggl

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

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

var ErrBadStatus = errors.New("bad http status")
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

func (t *Toggl) get(url string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err) // should not happen
	}

	req.SetBasicAuth(t.cfg.ApiToken, "api_token")
	req.Header.Add("Accept", "application/json")

	res, err := t.client.Do(req)

	if err != nil {
		return nil, err
	} else if res.StatusCode >= 400 {
		return nil, ErrBadStatus
	}

	return res.Body, nil
}

type timerResponse struct {
	Id          int
	Description string
	Duration    int64
	Start       time.Time
	End         time.Time
	ProjectId   int
	Tags        []string
}

func (t *Toggl) CurrentTimer() (*Timer, error) {
	body, err := t.get("https://api.track.toggl.com/api/v9/me/time_entries/current")

	if err != nil {
		return nil, err
	}

	defer body.Close()
	decoder := json.NewDecoder(body)

	var data timerResponse
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}

	timer := data.toTimer()
	if timer.Id == 0 {
		return nil, ErrNoTimer
	}

	return timer, nil
}
