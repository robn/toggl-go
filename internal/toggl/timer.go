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
	duration    int64
	projectId   int
}

func timerFromResponseBody(body io.Reader) (*Timer, error) {
	decoder := json.NewDecoder(body)

	// we decode here into a struct with all public members, then hide some of
	// them publicly
	var data struct {
		Id          int
		Description string
		Duration    int64
		Start       time.Time
		End         time.Time
		ProjectId   int `json:"project_id"`
		WorkspaceId int `json:"workspace_id"`
		Tags        []string
	}

	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}

	if data.Id == 0 {
		return nil, ErrNoTimer
	}

	return &Timer{
		Id:          data.Id,
		Description: data.Description,
		Tags:        data.Tags,
		WorkspaceId: data.WorkspaceId,
		duration:    data.Duration,
		projectId:   data.ProjectId,
	}, nil
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
	proj := "--" // TODO

	tagStr := ""
	if len(t.Tags) > 0 {
		prefixed := make([]string, 0, len(t.Tags))

		for _, tag := range t.Tags {
			prefixed = append(prefixed, "#"+tag)
		}

		tagStr = ", " + strings.Join(prefixed, ", ")
	}

	return fmt.Sprintf("%s (%s%s)", t.Description, proj, tagStr)
}
