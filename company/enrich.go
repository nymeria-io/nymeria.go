package company

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"git.nymeria.io/nymeria.go"
	"git.nymeria.io/nymeria.go/internal/api"
)

type EnrichParams struct {
	Website string
	Profile string
	Name    string
}

func (e EnrichParams) Invalid() bool {
	return len(e.Website) == 0 && len(e.Profile) == 0 && len(e.Name) == 0
}

func (e EnrichParams) URL() string {
	return fmt.Sprintf(
		"website=%s&profile=%s&name=%s",
		url.QueryEscape(e.Website),
		url.QueryEscape(e.Profile),
		url.QueryEscape(e.Name),
	)
}

func Enrich(params EnrichParams) (*Company, error) {
	log.Println("Enriching...")

	if params.Invalid() {
		return nil, nymeria.ErrInvalidParameters
	}

	req, err := api.Request("GET", fmt.Sprintf("/company/enrich?%s", params.URL()), nil)

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
		Status int     `json:"status"`
		Data   Company `json:"data"`
	}

	if err := json.Unmarshal(bs, &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}
