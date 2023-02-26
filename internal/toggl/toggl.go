package toggl

import (
	"errors"
	"fmt"
	"io"
	"log"
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

var BadStatus = errors.New("bad http status")

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
