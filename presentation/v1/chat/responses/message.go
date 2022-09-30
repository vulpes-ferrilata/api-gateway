package responses

type Message struct {
	ID     string `json:"id"`
	RoomID string `json:"roomID"`
	UserID string `json:"userID"`
	Detail string `json:"detail"`
}
