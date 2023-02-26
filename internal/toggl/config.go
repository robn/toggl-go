package toggl

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	ApiToken         string         `toml:"api_token"`
	ProjectShortcuts map[string]int `toml:"project_shortcuts"`
	projectsById     map[int]string
}

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

	byId := make(map[int]string)
	for name, id := range t.cfg.ProjectShortcuts {
		byId[id] = name
	}

	t.cfg.projectsById = byId

	return nil
}

// just a utility function which is only public so that it's easy to dump from
// one of the commands
func (t *Toggl) DebugConfig() {
	fmt.Fprintf(os.Stderr, "%#v\n", t.cfg)
}
