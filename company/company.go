package company

type Company struct {
	ID           string `json:"id"`
	Size         string `json:"size"`
	Name         string `json:"name"`
	Industry     string `json:"industry"`
	Founded      string `json:"founded"`
	WebsiteURL   string `json:"website_url"`
	LinkedinID   int    `json:"linkedin_id"`
	LinkedinName string `json:"linkedin_name"`
	TwitterName  string `json:"twitter_name"`
	FacebookName string `json:"facebook_name"`
	Location     string `json:"location"`
}
