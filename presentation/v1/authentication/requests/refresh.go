package requests

type Refresh struct {
	RefreshToken string `json:"refreshToken" validate:"required,jwt"`
}
