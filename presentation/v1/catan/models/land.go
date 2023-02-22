package models

type Land struct {
	ID       string `json:"id"`
	Q        int    `json:"q"`
	R        int    `json:"r"`
	Location string `json:"location"`
}
