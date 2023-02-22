package models

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required,jwt"`
}
