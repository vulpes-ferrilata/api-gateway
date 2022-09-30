package responses

type ResourceCard struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	IsSelected bool   `json:"isSelected"`
}
