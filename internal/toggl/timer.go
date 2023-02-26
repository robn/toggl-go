package toggl

import (
	"fmt"
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

func (r timerResponse) toTimer() *Timer {
	return &Timer{
		Id:          r.Id,
		Description: r.Description,
		Tags:        r.Tags,
		WorkspaceId: r.WorkspaceId,
		duration:    r.Duration,
		projectId:   r.ProjectId,
	}
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
