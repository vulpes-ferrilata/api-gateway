package models

type Message struct {
	ID     string `json:"id"`
	UserID string `json:"userID"`
	Detail string `json:"detail"`
}
