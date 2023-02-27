package toggl

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	ApiToken         string                       `toml:"api_token"`
	WorkspaceId      int                          `toml:"workspace_id"`
	ProjectShortcuts map[string]int               `toml:"project_shortcuts"`
	TaskShortcuts    map[string]map[string]string `toml:"task_shortcuts"`
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

	_, err = toml.Decode(string(tomlData), &t.Config)
	if err != nil {
		return fmt.Errorf("bad config file: %w", err)
	}

	byId := make(map[int]string)
	for name, id := range t.Config.ProjectShortcuts {
		byId[id] = name
	}

	t.Config.projectsById = byId

	return nil
}
