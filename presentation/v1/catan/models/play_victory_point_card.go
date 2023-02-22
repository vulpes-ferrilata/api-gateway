package models

type PlayVictoryPointCardRequest struct {
	DevelopmentCardID string `json:"developmentCardID" validate:"required,objectid"`
}
