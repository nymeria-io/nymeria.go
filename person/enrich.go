package person

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"git.nymeria.io/nymeria.go"
	"git.nymeria.io/nymeria.go/internal/api"
)

type EnrichParams struct {
	Profile string
	Email   string
	LID     string
	Filter  string
	Require string
}

func (e EnrichParams) Invalid() bool {
	return len(e.Profile) == 0 && len(e.Email) == 0 && len(e.LID) == 0
}

func (e EnrichParams) URL() string {
	return fmt.Sprintf(
		"profile=%s&email=%s&lid=%s&filter=%s&require=%s",
		url.QueryEscape(e.Profile),
		url.QueryEscape(e.Email),
		url.QueryEscape(e.LID),
		url.QueryEscape(e.Filter),
		url.QueryEscape(e.Require),
	)
}

func Enrich(params EnrichParams) (*Person, error) {
	if params.Invalid() {
		return nil, nymeria.ErrInvalidParameters
	}

	req, err := api.Request("GET", fmt.Sprintf("/person/enrich?%s", params.URL()), nil)

	if err != nil {
		return nil, err
	}

	resp, err := api.Client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		if e, ok := nymeria.ErrMap[resp.StatusCode]; ok {
			return nil, e
		}

		return nil, nymeria.ErrServerError
	}

	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var response struct {
		Status int    `json:"status"`
		Data   Person `json:"data"`
	}

	if err := json.Unmarshal(bs, &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}
