package models

type GamePagination struct {
	Total int     `json:"total"`
	Data  []*Game `json:"data"`
}
