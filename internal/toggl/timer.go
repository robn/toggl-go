package toggl

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

type Timer struct {
	Id          int
	Description string
	Tags        []string
	WorkspaceId int
	Project     string
	Start       time.Time
	End         time.Time
	duration    int64
	projectId   int
}

// We decode into a struct with all public members, then hide some of them publicly
type timerData struct {
	Id          int
	Description string
	Duration    int64
	Start       time.Time
	End         time.Time
	ProjectId   int `json:"project_id"`
	WorkspaceId int `json:"workspace_id"`
	Tags        []string
}

func (t *Toggl) timerFromData(data timerData) *Timer {
	project, ok := t.Config.projectsById[data.ProjectId]
	if !ok {
		project = "--"
	}

	return &Timer{
		Id:          data.Id,
		Description: data.Description,
		Tags:        data.Tags,
		WorkspaceId: data.WorkspaceId,
		Project:     project,
		Start:       data.Start,
		End:         data.End,
		duration:    data.Duration,
		projectId:   data.ProjectId,
	}
}

func (t *Toggl) timerFromResponseBody(body io.Reader) (*Timer, error) {
	var data timerData
	decoder := json.NewDecoder(body)

	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}

	if data.Id == 0 {
		return nil, ErrNoTimer
	}

	return t.timerFromData(data), nil
}

func (t *Toggl) timersFromResponseBody(body io.Reader) ([]*Timer, error) {
	var data []timerData
	decoder := json.NewDecoder(body)

	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}

	ret := make([]*Timer, 0, len(data))
	for _, td := range data {
		ret = append(ret, t.timerFromData(td))
	}

	return ret, nil
}

func (t Timer) Duration() time.Duration {
	dur := t.duration

	if dur < 0 {
		dur = time.Now().Unix() + int64(dur)
	}

	// duration is in nanoseconds, and we have seconds
	return time.Duration(dur * 1e9)
}

func (t Timer) OnelineDesc() string {
	tagStr := ""
	if len(t.Tags) > 0 {
		prefixed := make([]string, 0, len(t.Tags))

		for _, tag := range t.Tags {
			prefixed = append(prefixed, "#"+tag)
		}

		tagStr = ", " + strings.Join(prefixed, ", ")
	}

	return fmt.Sprintf("%s (%s%s)", t.Description, t.Project, tagStr)
}
