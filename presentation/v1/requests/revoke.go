package requests

type Revoke struct {
	RefreshToken string `json:"refreshToken" validate:"required,jwt"`
}
