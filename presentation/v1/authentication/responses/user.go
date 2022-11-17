package responses

type User struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName,omitempty"`
}
