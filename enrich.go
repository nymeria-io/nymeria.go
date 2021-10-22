package nymeria

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
)

// Enrichment is a response from the enrichment API. The results of the enrich
// call, if successful contains a person's information based on one or more of
// the query parameters.
type Enrichment struct {
	Status           string `json:"status"`
	DeveloperMessage string `json:"developer_message,omitempty"`

	Meta struct {
		URL        string `json:"url,omitempty"`
		Identifier string `json:"identifier,omitempty"`
	} `json:"meta"`

	Usage struct {
		Used  int `json:"used"`
		Limit int `json:"limit"`
	} `json:"usage"`

	Data struct {
		Bio struct {
			FirstName      string `json:"first_name,omitempty"`
			LastName       string `json:"last_name,omitempty"`
			Company        string `json:"company,omitempty"`
			CompanyWebsite string `json:"company_website,omitempty"`
		} `json:"bio,omitempty"`

		Emails []struct {
			Type    string `json:"type,omitempty"`
			Name    string `json:"name,omitempty"`
			Domain  string `json:"domain,omitempty"`
			Address string `json:"address,omitempty"`
		} `json:"emails,omitempty"`

		PhoneNumbers []struct {
			Number string `json:"number,omitempty"`
		} `json:"phone_numbers,omitempty"`

		Social []struct {
			Type string `json:"type,omitempty"`
			ID   string `json:"id,omitempty"`
			URL  string `json:"url,omitempty"`
		} `json:"social,omitempty"`
	} `json:"data"`
}

// Enrich takes a url as a first parameter.
func Enrich(u string) (*Enrichment, error) {
	req, err := request("GET", fmt.Sprintf("/enrich?url=%s", url.QueryEscape(u)))

	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var response Enrichment

	if err := json.Unmarshal(bs, &response); err != nil {
		return nil, err
	}

	if response.Status != "success" {
		return nil, ErrInvalidRequest
	}

	return &response, nil
}
