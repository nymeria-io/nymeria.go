package nymeria

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
)

var (
	// ErrInvalidRequest is returned any time a request looks malformed or invalid.
	// This can be due to a bad parameter, encoding, authentication, etc.
	ErrInvalidRequest = errors.New(`error: the request failed; perhaps you haven't authenticated or the request was malformed`)
)

const (
	// ProfessionalEmailFilter is a valid EnrichParams Filter. If specified, no professional
	// emails will be returned. As a result, only personal emails will be returned.
	ProfessionalEmailFilter = "professional-emails"
)

// EnrichParams are the parameters used for the enrichment look up. One or more
// can be supplied to Enrich.
type EnrichParams struct {
	URL        string                 `json:"url"`
	Email      string                 `json:"email"`
	Identifier string                 `json:"identifier"`
	Filter     string                 `json:"filter"`
	Require    string                 `json:"require"`
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

// PeopleQuery is the query parameters for the people search API.
type PeopleQuery struct {
	Start    int      `json:"start,omitempty"`
	Q        string   `json:"q,omitempty"`
	Location string   `json:"location,omitempty"`
	Country  string   `json:"country,omitempty"`
	Title    string   `json:"title,omitempty"`
	Company  string   `json:"company,omitempty"`
	Skills   []string `json:"skills,omitempty"`
	HasEmail bool     `json:"has_email,omitempty"`
	HasPhone bool     `json:"has_phone,omitempty"`
}

// Params will convert the query into a query string.
func (p PeopleQuery) Params() string {
	var sb strings.Builder

	if len(p.Q) > 0 {
		sb.WriteString(fmt.Sprintf(`&q=%s`, url.QueryEscape(p.Q)))
	}

	if len(p.Location) > 0 {
		sb.WriteString(fmt.Sprintf(`&location=%s`, url.QueryEscape(p.Location)))
	}

	if len(p.Country) > 0 {
		sb.WriteString(fmt.Sprintf(`&country=%s`, url.QueryEscape(p.Country)))
	}

	if len(p.Title) > 0 {
		sb.WriteString(fmt.Sprintf(`&title=%s`, url.QueryEscape(p.Title)))
	}

	if len(p.Company) > 0 {
		sb.WriteString(fmt.Sprintf(`&company=%s`, url.QueryEscape(p.Company)))
	}

	if p.HasEmail {
		sb.WriteString(`&has_email=true`)
	}

	if p.HasPhone {
		sb.WriteString(`&has_phone=true`)
	}

	if len(p.Skills) > 0 {
		sb.WriteString(fmt.Sprintf(`&skills=%s`, url.QueryEscape(strings.Join(p.Skills, ","))))
	}

	return sb.String()
}

// PeoplePreview is a preview of available data for each person returned by
// the search API.
type PeoplePreview struct {
	UUID          string   `json:"uuid"`
	FirstName     string   `json:"first_name,omitempty"`
	LastName      string   `json:"last_name,omitempty"`
	Title         string   `json:"title,omitempty"`
	Company       string   `json:"company,omitempty"`
	Location      string   `json:"location,omitempty"`
	Country       string   `json:"country,omitempty"`
	AvailableData []string `json:"available_data,omitempty"`
}

// PeopleResponse is the raw response returned by the API.
type PeopleResponse struct {
	Status string `json:"status"`

	Meta struct {
		Query struct {
			Rows     int      `json:"rows"`
			Start    int      `json:"start"`
			Q        string   `json:"q"`
			Location string   `json:"location"`
			Company  string   `json:"company"`
			Title    string   `json:"title"`
			HasEmail bool     `json:"has_email"`
			HasPhone bool     `json:"has_phone"`
			Skills   []string `json:"skills"`
		} `json:"query"`
		Results struct {
			Total int `json:"total"`
		} `json:"results"`
	} `json:"meta"`

	Data []PeoplePreview `json:"data"`
}

// RevealedPeople are people that has been revealed (contact and other data unlocked).
type RevealedPeople struct {
	Status string `json:"status"`

	Meta struct {
		UUIDS []string `json:"uuids"`
	} `json:"meta"`

	Usage struct {
		Used  int `json:"used"`
		Limit int `json:"limit"`
	} `json:"usage"`

	Data []Person `json:"data"`
}

// People will perform a search query for people and return a preview of the people.
func People(q *PeopleQuery) (*PeopleResponse, error) {
	req, err := request("GET", fmt.Sprintf("/people?1=1%s", q.Params()), nil)

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

	var response PeopleResponse

	if err := json.Unmarshal(bs, &response); err != nil {
		return nil, err
	}

	if response.Status != "success" {
		return nil, ErrInvalidRequest
	}

	return &response, nil
}

// RevealPeople will reveal data for a person given their uuid. Each reveal
// will consume one credit.
func RevealPeople(uuids []string) (*RevealedPeople, error) {
	type params struct {
		UUIDs []string `json:"uuids"`
	}

	bs, err := json.Marshal(params{UUIDs: uuids})

	if err != nil {
		return nil, err
	}

	req, err := request("POST", "/people", bytes.NewBuffer(bs))

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

	var response RevealedPeople

	if err := json.Unmarshal(bs, &response); err != nil {
		return nil, err
	}

	if response.Status != "success" {
		return nil, ErrInvalidRequest
	}

	return &response, nil
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
