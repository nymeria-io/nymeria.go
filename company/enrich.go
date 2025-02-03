package company

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/nymeriaio/nymeria.go"
	"github.com/nymeriaio/nymeria.go/internal/api"
)

type EnrichParams struct {
	Website    string
	Profile    string
	Name       string
	LinkedinID int
}

func (e EnrichParams) Invalid() bool {
	return len(e.Website) == 0 && len(e.Profile) == 0 && len(e.Name) == 0 && e.LinkedinID == 0
}

func (e EnrichParams) URL() string {
	var query strings.Builder

	query.WriteString("1=1")

	if len(e.Website) > 0 {
		query.WriteString(fmt.Sprintf("&website=%s", url.QueryEscape(e.Website)))
	}

	if len(e.Name) > 0 {
		query.WriteString(fmt.Sprintf("&name=%s", url.QueryEscape(e.Name)))
	}

	if len(e.Profile) > 0 {
		query.WriteString(fmt.Sprintf("&profile=%s", url.QueryEscape(e.Profile)))
	}

	if e.LinkedinID > 0 {
		query.WriteString(fmt.Sprintf("&linkedin_id=%d", e.LinkedinID))
	}

	return query.String()
}

func Enrich(params EnrichParams) (*Company, error) {
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

	bs, err := io.ReadAll(resp.Body)

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
