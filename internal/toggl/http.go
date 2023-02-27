package toggl

import (
	"net/http"
	"os"
)

func (t *Toggl) doRequest(req *http.Request) (*http.Response, error) {
	// default headers
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(t.Config.ApiToken, "api_token")

	if req.Body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	res, err := t.client.Do(req)

	if err != nil {
		return nil, err
	} else if res.StatusCode >= 400 {
		// just bail, I'm never going to do anything useful with this.
		dumpResponseAndExit(res)
	}

	return res, nil
}

func (t *Toggl) get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err) // should not happen
	}

	return t.doRequest(req)
}

func dumpResponseAndExit(res *http.Response) {
	err := res.Write(os.Stderr)
	if err != nil {
		panic("could not write out response to stderr")
	}

	os.Exit(1)
}
