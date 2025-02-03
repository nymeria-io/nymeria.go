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

type SearchParams struct {
	Name     string
	Location string
	Country  string
	Industry string
	Size     string
	Limit    int /* how many records to retrieve starting (default: 0) */
	Offset   int /* from which record to start */
}

func (s SearchParams) Invalid() bool {
	return len(s.Name) == 0 && len(s.Location) == 0 && len(s.Country) == 0 && len(s.Industry) == 0 && len(s.Size) == 0
}

func (s SearchParams) URL() string {
	if s.Limit <= 0 || s.Limit > 100 {
		s.Limit = 10
	}

	var query strings.Builder

	query.WriteString(fmt.Sprintf("limit=%d", s.Limit))
	query.WriteString(fmt.Sprintf("&offset=%d", s.Offset))

	if len(s.Name) > 0 {
		query.WriteString(fmt.Sprintf("&name=%s", url.QueryEscape(s.Name)))
	}

	if len(s.Size) > 0 {
		query.WriteString(fmt.Sprintf("&size=%s", url.QueryEscape(s.Size)))
	}

	if len(s.Location) > 0 {
		query.WriteString(fmt.Sprintf("&location=%s", url.QueryEscape(s.Location)))
	}

	if len(s.Country) > 0 {
		query.WriteString(fmt.Sprintf("&country=%s", url.QueryEscape(s.Country)))
	}

	if len(s.Industry) > 0 {
		query.WriteString(fmt.Sprintf("&industry=%s", url.QueryEscape(s.Industry)))
	}

	return query.String()
}

func Search(params SearchParams) ([]Company, error) {
	if params.Invalid() {
		return nil, nymeria.ErrInvalidParameters
	}

	req, err := api.Request("GET", fmt.Sprintf("/company/search?%s", params.URL()), nil)

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
		Status int       `json:"status"`
		Data   []Company `json:"data"`
	}

	if err := json.Unmarshal(bs, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}
