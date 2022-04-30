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

// EnrichParams are the parameters used for the enrichment look up. One or more
// can be supplied to Enrich.
type EnrichParams struct {
	URL        string                 `json:"url"`
	Email      string                 `json:"email"`
	Identifier string                 `json:"identifier"`
	Custom     map[string]interface{} `json:"custom,omitempty"`
}

// Enrichment is a response from the enrichment API. The results of the enrich
// call, if successful contains a person's information based on one or more of
// the query parameters.
type Enrichment struct {
	Status           string `json:"status"`
	DeveloperMessage string `json:"developer_message,omitempty"`

	Meta EnrichParams `json:"meta"`

	Usage struct {
		Used  int `json:"used"`
		Limit int `json:"limit"`
	} `json:"usage"`

	Data Person `json:"data"`
}

// Person is a successful enrichment.
type Person struct {
	Bio struct {
		FirstName      string `json:"first_name,omitempty"`
		LastName       string `json:"last_name,omitempty"`
		Title          string `json:"title,omitempty"`
		Location       string `json:"location,omitempty"`
		Country        string `json:"country,omitempty"`
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

// Enrich takes one or more EnrichParams and will return an enrichment result
// for each enrichment.
func Enrich(args ...EnrichParams) ([]Enrichment, error) {
	var es []Enrichment

	switch len(args) {
	case 0: /* error */
		return es, ErrInvalidRequest
	case 1: /* single enrichment*/
		enrichment, err := enrich(args[0])

		if err != nil {
			return es, err
		}

		es = append(es, *enrichment)

		return es, nil
	default: /* bulk enrichment*/
		es, err := bulk(args...)

		if err != nil {
			return es, err
		}

		return es, nil
	}

	return es, ErrInvalidRequest
}

// Performs a single enrichment.
func enrich(params EnrichParams) (*Enrichment, error) {
	bs, err := json.Marshal(params)

	if err != nil {
		return nil, err
	}

	req, err := request("POST", "/enrich", bytes.NewBuffer(bs))

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

	var response Enrichment

	if err := json.Unmarshal(bs, &response); err != nil {
		return nil, err
	}

	if response.Status != "success" {
		return nil, ErrInvalidRequest
	}

	return &response, nil
}

// Performs a bulk enrichment and converts the results into a slice of
// enrichments.
func bulk(params ...EnrichParams) ([]Enrichment, error) {
	var es []Enrichment

	if len(params) == 0 {
		return es, ErrInvalidRequest
	}

	type BulkEnrichment struct {
		Status           string `json:"status"`
		DeveloperMessage string `json:"developer_message,omitempty"`

		Usage struct {
			Used  int `json:"used"`
			Limit int `json:"limit"`
		} `json:"usage"`

		Data []struct {
			Meta   EnrichParams `json:"meta"`
			Result Person       `json:"result"`
		} `json:"data"`
	}

	type BulkPayload struct {
		People []EnrichParams `json:"people"`
	}

	var payload BulkPayload

	payload.People = append(payload.People, params...)

	bs, err := json.Marshal(payload)

	if err != nil {
		return es, err
	}

	req, err := request("POST", "/bulk-enrich", bytes.NewBuffer(bs))

	if err != nil {
		return es, err
	}

	resp, err := client.Do(req)

	if err != nil {
		return es, err
	}

	defer resp.Body.Close()

	bs, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		return es, err
	}

	var response BulkEnrichment

	if err := json.Unmarshal(bs, &response); err != nil {
		return es, err
	}

	if response.Status != "success" {
		return es, ErrInvalidRequest
	}

	for _, result := range response.Data {
		es = append(es, Enrichment{
			Status:           response.Status,
			DeveloperMessage: response.DeveloperMessage,
			Usage:            response.Usage,
			Meta:             result.Meta,
			Data:             result.Result,
		})
	}

	return es, nil
}
