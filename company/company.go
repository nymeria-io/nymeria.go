package company

type Company struct {
	ID                string `json:"id"`
	EmployeeCountFrom int    `json:"employee_count_from"`
	EmployeeCountTo   int    `json:"employee_count_to"`
	Size              string `json:"size"`
	Name              string `json:"name"`
	Industry          string `json:"industry"`
	Founded           int    `json:"founded"`
	WebsiteURL        string `json:"website_url"`
	LinkedinID        string `json:"linkedin_id"`
	LinkedinName      string `json:"linkedin_name"`
	TwitterName       string `json:"twitter_name"`
	FacebookName      string `json:"facebook_name"`
	LocationCountry   string `json:"location_country"`
	LocationGeo       string `json:"location_geo"`
	LocationContinent string `json:"location_continent"`
	LocationName      string `json:"location_name"`
	LocationRegion    string `json:"location_region"`
	UpdatedAt         string `json:"updated_at"`
}
