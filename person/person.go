package person

type Person struct {
	ID                    string         `json:"id"`
	FirstName             *string        `json:"first_name"`
	LastName              *string        `json:"last_name"`
	FullName              *string        `json:"full_name"`
	Gender                *string        `json:"gender"`
	Age                   *int           `json:"age"`
	BirthYear             *string        `json:"birth_year"`
	BirthDate             *string        `json:"birth_date"`
	WorkEmail             *string        `json:"work_email"`
	PersonalEmails        []string       `json:"personal_emails"`
	Emails                []EmailAddress `json:"emails"`
	MobilePhone           *string        `json:"mobile_phone"`
	PhoneNumbers          []string       `json:"phone_numbers"`
	Industry              *string        `json:"industry"`
	LocationName          *string        `json:"location_name"`
	LocationLastUpdated   *string        `json:"location_last_updated"`
	LocationCountry       *string        `json:"location_country"`
	InferredExperience    *int           `json:"inferred_years_of_experience"`
	InferredSalary        *string        `json:"inferred_salary"`
	JobTitle              *string        `json:"job_title"`
	JobTitleRole          *string        `json:"job_title_role"`
	JobTitleLevels        []string       `json:"job_title_levels"`
	JobStartDate          *string        `json:"job_start_date"`
	JobCompanyName        *string        `json:"job_company_name"`
	JobCompanyURL         *string        `json:"job_company_website"`
	JobCompanyFounded     *string        `json:"job_company_founded"`
	JobCompanySize        *string        `json:"job_company_size"`
	JobCompanyLinkedinURL *string        `json:"job_company_linkedin_url"`
	JobLastUpdated        *string        `json:"job_last_updated"`
	JobSummary            *string        `json:"job_summary"`
	Skills                []string       `json:"skills"`
	Interests             []string       `json:"interests"`
	LinkedinUsername      *string        `json:"linkedin_username"`
	LinkedinURL           *string        `json:"linkedin_url"`
	LinkedinID            *string        `json:"linkedin_id"`
	LinkedinConnections   *int           `json:"linkedin_connections"`
	FacebookUsername      *string        `json:"facebook_username"`
	FacebookURL           *string        `json:"facebook_url"`
	FacebookID            *string        `json:"facebook_id"`
	TwitterUsername       *string        `json:"twitter_username"`
	TwitterURL            *string        `json:"twitter_url"`
	GithubUsername        *string        `json:"github_username"`
	GithubURL             *string        `json:"github_url"`
	Profiles              []SocialLink   `json:"profiles"`
	LinkedinSummary       *string        `json:"linkedin_summary"`
	Education             []Education    `json:"education"`
	Experience            []Experience   `json:"experience"`
	Certificates          []Certificate  `json:"certificates"`
	Languages             []Language     `json:"languages"`
}

type Language struct {
	Name        string `json:"name"`
	Proficiency int    `json:"proficiency"`
}

type EmailAddress struct {
	Type   string `json:"type"` /* personal, professional, educational, disposable */
	Name   string `json:"name"`
	Domain string `json:"domain"`
	Full   string `json:"address"`
}

type SocialLink struct {
	Network  string `json:"network"`
	URL      string `json:"url"`
	Username string `json:"username"`
}

type Education struct {
	Majors    []string `json:"majors"`
	EndDate   *string  `json:"end_date"`
	StartDate *string  `json:"start_date"`
	School    *struct {
		ID          *string `json:"id"`
		Name        *string `json:"name"`
		Type        *string `json:"type"`
		Domain      *string `json:"domain"`
		Website     *string `json:"website"`
		LinkedinID  *string `json:"linkedin_id"`
		LinkedinURL *string `json:"linkedin_url"`
	} `json:"school,omitempty"`
}

type Experience struct {
	Title *struct {
		Name    *string  `json:"name"`
		Role    *string  `json:"role"`
		SubRole *string  `json:"sub_role"`
		Levels  []string `json:"levels"`
	} `json:"title,omitempty"`

	EndDate   *string `json:"end_date"`
	StartDate *string `json:"start_date"`
}

type Certificate struct {
	Name         *string `json:"name"`
	Organization *string `json:"organization"`
	EndDate      *string `json:"end_date"`
	StartDate    *string `json:"start_date"`
}
