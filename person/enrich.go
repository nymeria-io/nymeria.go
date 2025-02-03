package person

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/nymeriaio/nymeria.go"
	"github.com/nymeriaio/nymeria.go/internal/api"
)

type BulkEnrichParams struct {
	Params   EnrichParams `json:"params"`
	MetaData interface{}  `json:"metadata"`
}

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
	var query strings.Builder

	prefix := ""

	if len(e.Profile) > 0 {
		query.WriteString(fmt.Sprintf("profile=%s", url.QueryEscape(e.Profile)))
		prefix = "&"
	}

	if len(e.Email) > 0 {
		query.WriteString(fmt.Sprintf("%semail=%s", prefix, url.QueryEscape(e.Email)))
		prefix = "&"
	}

	if len(e.LID) > 0 {
		query.WriteString(fmt.Sprintf("%slid=%s", prefix, url.QueryEscape(e.LID)))
		prefix = "&"
	}

	if len(e.Filter) > 0 {
		query.WriteString(fmt.Sprintf("%sfilter=%s", prefix, url.QueryEscape(e.Filter)))
		prefix = "&"
	}

	if len(e.Require) > 0 {
		query.WriteString(fmt.Sprintf("%srequire=%s", prefix, url.QueryEscape(e.Require)))
		prefix = "&"
	}

	return query.String()
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

	bs, err := io.ReadAll(resp.Body)

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

func BulkEnrich(params ...BulkEnrichParams) ([]Person, error) {
	if len(params) == 0 {
		return nil, nymeria.ErrInvalidParameters
	}

	bs, err := json.Marshal(map[string]interface{}{
		"requests": params,
	})

	if err != nil {
		return nil, err
	}

	req, err := api.Request("POST", "/person/enrich/bulk", bytes.NewBuffer(bs))

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

	bs, err = io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var response []struct {
		Status int    `json:"status"`
		Data   Person `json:"data"`
	}

	if err := json.Unmarshal(bs, &response); err != nil {
		return nil, err
	}

	var records []Person

	for _, v := range response {
		records = append(records, v.Data)
	}

	return records, nil
}
