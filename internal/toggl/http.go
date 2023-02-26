package toggl

import (
	"errors"
	"io"
	"net/http"
	"os"
)

var ErrBadStatus = errors.New("bad http status")

func (t *Toggl) doRequest(req *http.Request) (*http.Response, error) {
	// default headers
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(t.cfg.ApiToken, "api_token")

	res, err := t.client.Do(req)

	if err != nil {
		return nil, err
	} else if res.StatusCode >= 400 {
		return nil, ErrBadStatus
	}

	return res, nil
}

func (t *Toggl) get(url string) (io.ReadCloser, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err) // should not happen
	}

	res, err := t.doRequest(req)
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}

func dumpResponseAndExit(res *http.Response) {
	err := res.Write(os.Stderr)
	if err != nil {
		panic("could not write out response to stderr")
	}

	os.Exit(1)
}
