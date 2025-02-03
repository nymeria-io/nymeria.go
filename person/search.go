package person

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/nymeria-io/nymeria.go"
)

type SearchParams struct {
	FirstName string
	LastName  string
	Title     string
	Company   string
	Country   string
	Location  string
	Industry  string
	Limit     int /* how many records to retrieve starting (default: 0) */
	Offset    int /* from which record to start */
}

func (s SearchParams) Invalid() bool {
	return len(s.FirstName) == 0 && len(s.LastName) == 0 && len(s.Title) == 0 && len(s.Company) == 0 && len(s.Country) == 0 && len(s.Industry) == 0
}

func (s SearchParams) URL() string {
	if s.Limit <= 0 || s.Limit > 100 {
		s.Limit = 10
	}

	var query strings.Builder

	query.WriteString(fmt.Sprintf("limit=%d", s.Limit))
	query.WriteString(fmt.Sprintf("&offset=%d", s.Offset))

	if len(s.FirstName) > 0 {
		query.WriteString(fmt.Sprintf("&first_name=%s", url.QueryEscape(s.FirstName)))
	}

	if len(s.LastName) > 0 {
		query.WriteString(fmt.Sprintf("&last_name=%s", url.QueryEscape(s.LastName)))
	}

	if len(s.Title) > 0 {
		query.WriteString(fmt.Sprintf("&title=%s", url.QueryEscape(s.Title)))
	}

	if len(s.Company) > 0 {
		query.WriteString(fmt.Sprintf("&company=%s", url.QueryEscape(s.Company)))
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

func Search(params SearchParams) ([]Person, error) {
	if params.Invalid() {
		return nil, nymeria.ErrInvalidParameters
	}

	req, err := nymeria.Request("GET", fmt.Sprintf("/person/search?%s", params.URL()), nil)

	if err != nil {
		return nil, err
	}

	resp, err := nymeria.Client.Do(req)

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
		Data     []Person    `json:"data"`
		Status   int         `json:"status"`
		MetaData interface{} `json:"metadata"`
		Total    int         `json:"total"`
	}

	fmt.Println(string(bs))

	if err := json.Unmarshal(bs, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}
