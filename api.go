package nymeria

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
)

var (
	// ErrInvalidRequest is returned any time a request looks malformed or invalid.
	// This can be due to a bad parameter, encoding, authentication, etc.
	ErrInvalidRequest = errors.New(`error: the request looks invalid or malformed`)
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

// BulkEnrichment is an aggregate response of Enrichments.
type BulkEnrichment struct {
	Status           string `json:"status"`
	DeveloperMessage string `json:"developer_message,omitempty"`

	Usage struct {
		Used  int `json:"used"`
		Limit int `json:"limit"`
	} `json:"usage"`

	Data []struct {
		Meta struct {
			Custom     map[string]interface{} `json:"custom,omitempty"`
			URL        string                 `json:"url,omitempty"`
			Identifier string                 `json:"identifier,omitempty"`
		} `json:"meta"`

		Result struct {
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
		} `json:"result"`
	} `json:"data"`
}

// Person represents a single person query.
type Person struct {
	Custom     map[string]interface{} `json:"custom,omitempty"`
	URL        string                 `json:"url,omitempty"`
	Identifier string                 `json:"identifier,omitempty"`
}

// BulkPayload is the payload of urls and other query parameters.
type BulkPayload struct {
	People []Person `json:"people"`
}

//
// Verification is the result of an email verify request. The "result" can either be
// valid or invalid. The tags can include things like:
//
//		has_dns_mx, smtp_connectable, accepts_all, etc
//
type Verification struct {
	Status string `json:"status"`

	Meta struct {
		Email string `json:"email"`
	} `json:"meta"`

	Usage struct {
		Used  int `json:"used"`
		Limit int `json:"limit"`
	} `json:"usage"`

	Data struct {
		Result string   `json:"result"`
		Tags   []string `json:"tags"`
	} `json:"data"`
}

// Verify takes a professional email address and tries to verify its validity.
func Verify(email string) (*Verification, error) {
	req, err := request("GET", fmt.Sprintf("/verify?email=%s", url.QueryEscape(email)), nil)

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

	var response Verification

	if err := json.Unmarshal(bs, &response); err != nil {
		return nil, err
	}

	if response.Status != "success" {
		return nil, ErrInvalidRequest
	}

	return &response, nil
}

// Enrich takes a url as a first parameter.
func Enrich(u string) (*Enrichment, error) {
	req, err := request("GET", fmt.Sprintf("/enrich?url=%s", url.QueryEscape(u)), nil)

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

// BulkEnrich takes one or more custom queries and returns zero or more matches.
func BulkEnrich(u ...string) (*BulkEnrichment, error) {
	var payload BulkPayload

	for _, iu := range u {
		payload.People = append(payload.People, Person{
			URL: iu,
		})
	}

	bs, err := json.Marshal(payload)

	if err != nil {
		return nil, err
	}

	req, err := request("POST", "/bulk-enrich", bytes.NewBuffer(bs))

	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bs, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var response BulkEnrichment

	if err := json.Unmarshal(bs, &response); err != nil {
		return nil, err
	}

	if response.Status != "success" {
		return nil, ErrInvalidRequest
	}

	return &response, nil
}
