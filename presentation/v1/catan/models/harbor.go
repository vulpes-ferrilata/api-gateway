package models

type Harbor struct {
	ID   string `json:"id"`
	Q    int    `json:"q"`
	R    int    `json:"r"`
	Type string `json:"type"`
}
