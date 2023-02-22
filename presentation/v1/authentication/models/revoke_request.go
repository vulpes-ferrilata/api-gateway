package models

type RevokeRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required,jwt"`
}
