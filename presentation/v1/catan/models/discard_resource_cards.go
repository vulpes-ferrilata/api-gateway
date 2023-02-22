package models

type DiscardResourceCardsRequest struct {
	ResourceCardIDs []string `json:"resourceCardIDs" validate:"required,unique"`
}
