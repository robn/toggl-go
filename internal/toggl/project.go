package toggl

import (
	"encoding/json"
	"io"
)

type Project struct {
	Id          int
	Name        string
}

func (t *Toggl) projectsFromResponseBody(body io.Reader) (map[int]Project, error) {
	var projects []Project
	decoder := json.NewDecoder(body)

	if err := decoder.Decode(&projects); err != nil {
		return nil, err
	}

	ret := make(map[int]Project, len(projects))
	for _, p := range projects {
		ret[p.Id] = p
	}

	return ret, nil
}
