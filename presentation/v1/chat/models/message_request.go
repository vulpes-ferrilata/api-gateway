package models

type CreateMessageRequest struct {
	RoomID string `json:"roomID"`
	Detail string `json:"detail"`
}
