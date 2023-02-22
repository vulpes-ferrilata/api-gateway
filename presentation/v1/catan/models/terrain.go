package models

type Terrain struct {
	ID     string  `json:"id"`
	Q      int     `json:"q"`
	R      int     `json:"r"`
	Number int     `json:"number"`
	Type   string  `json:"type"`
	Harbor *Harbor `json:"harbor"`
	Robber *Robber `json:"robber"`
}
