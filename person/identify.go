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

type IdentifyParams struct {
	FirstName string
	LastName  string
	Name      string
	Location  string
	Country   string
	Filter    string
	Require   string
}

func (i IdentifyParams) Invalid() bool {
	return len(i.FirstName) == 0 && len(i.LastName) == 0 && len(i.Name) == 0 && len(i.Location) == 0 && len(i.Country) == 0
}

func (i IdentifyParams) URL() string {
	return fmt.Sprintf(
		"first_name=%s&last_name=%s&name=%s&location=%s&country=%s&filter=%s&require=%s",
		url.QueryEscape(i.FirstName),
		url.QueryEscape(i.LastName),
		url.QueryEscape(i.Name),
		url.QueryEscape(i.Location),
		url.QueryEscape(i.Country),
		url.QueryEscape(i.Filter),
		url.QueryEscape(i.Require),
	)
}

func Identify(params IdentifyParams) ([]Person, error) {
	if params.Invalid() {
		return nil, nymeria.ErrInvalidParameters
	}

	req, err := api.Request("GET", fmt.Sprintf("/person/identify?%s", params.URL()), nil)

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

	var response []struct {
		MatchedOn []string `json:"matched_on"`
		Data      Person   `json:"data"`
	}

	if err := json.Unmarshal(bs, &response); err != nil {
		return nil, err
	}

	people := []Person{}

	for _, r := range response {
		people = append(people, r.Data)
	}

	return people, nil
}
