package email

type Verification struct {
	Result              string   `json:"result"`
	Flags               []string `json:"flags"`
	SuggestedCorrection string   `json:"suggested_correction"`
	ExecutionTime       int      `json:"execution_time"`
}
