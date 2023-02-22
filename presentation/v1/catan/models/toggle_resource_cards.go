package models

type ToggleResourceCardsRequest struct {
	ResourceCardIDs []string `json:"resourceCardIDs" validate:"required,unique"`
}
