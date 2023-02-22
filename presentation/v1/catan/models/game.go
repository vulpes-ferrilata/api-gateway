package models

type Game struct {
	ID             string `json:"id"`
	PlayerQuantity int    `json:"playerQuantity"`
	Status         string `json:"status"`
}
