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

type PersonPreview struct {
	ID             string `json:"id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	FullName       string `json:"full_name"`
	JobTitle       string `json:"job_title"`
	LocationName   string `json:"location_name"`
	JobCompanyName string `json:"job_company_name"`

	Gender                bool `json:"gender"`
	Age                   bool `json:"age"`
	BirthYear             bool `json:"birth_year"`
	BirthDate             bool `json:"birth_date"`
	WorkEmail             bool `json:"work_email"`
	PersonalEmails        bool `json:"personal_emails"`
	Emails                bool `json:"emails"`
	MobilePhone           bool `json:"mobile_phone"`
	PhoneNumbers          bool `json:"phone_numbers"`
	Industry              bool `json:"industry"`
	LocationLastUpdated   bool `json:"location_last_updated"`
	LocationCountry       bool `json:"location_country"`
	InferredExperience    bool `json:"inferred_years_of_experience"`
	InferredSalary        bool `json:"inferred_salary"`
	JobTitleRole          bool `json:"job_title_role"`
	JobTitleLevels        bool `json:"job_title_levels"`
	JobStartDate          bool `json:"job_start_date"`
	JobCompanyURL         bool `json:"job_company_website"`
	JobCompanyFounded     bool `json:"job_company_founded"`
	JobCompanySize        bool `json:"job_company_size"`
	JobCompanyLinkedinURL bool `json:"job_company_linkedin_url"`
	JobLastUpdated        bool `json:"job_last_updated"`
	JobSummary            bool `json:"job_summary"`
	Skills                bool `json:"skills"`
	Interests             bool `json:"interests"`
	LinkedinUsername      bool `json:"linkedin_username"`
	LinkedinURL           bool `json:"linkedin_url"`
	LinkedinID            bool `json:"linkedin_id"`
	LinkedinConnections   bool `json:"linkedin_connections"`
	FacebookUsername      bool `json:"facebook_username"`
	FacebookURL           bool `json:"facebook_url"`
	FacebookID            bool `json:"facebook_id"`
	TwitterUsername       bool `json:"twitter_username"`
	TwitterURL            bool `json:"twitter_url"`
	GithubUsername        bool `json:"github_username"`
	GithubURL             bool `json:"github_url"`
	Profiles              bool `json:"profiles"`
	LinkedinSummary       bool `json:"linkedin_summary"`
	Education             bool `json:"education"`
	Experience            bool `json:"experience"`
	Certificates          bool `json:"certificates"`
	Languages             bool `json:"languages"`
}

type PreviewParams struct {
	Profile string
	Email   string
	LID     string
	Filter  string
	Require string
}

func (e PreviewParams) Invalid() bool {
	return len(e.Profile) == 0 && len(e.Email) == 0 && len(e.LID) == 0
}

func (e PreviewParams) URL() string {
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

func Preview(params PreviewParams) (*PersonPreview, error) {
	if params.Invalid() {
		return nil, nymeria.ErrInvalidParameters
	}

	req, err := nymeria.Request("GET", fmt.Sprintf("/person/enrich/preview?%s", params.URL()), nil)

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
		Status int           `json:"status"`
		Data   PersonPreview `json:"data"`
	}

	if err := json.Unmarshal(bs, &response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}
