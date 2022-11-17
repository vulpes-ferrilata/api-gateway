package responses

type ResourceCard struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Offering bool   `json:"offering"`
}
