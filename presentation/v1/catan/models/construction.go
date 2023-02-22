package models

type Construction struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Land *Land  `json:"land"`
}
