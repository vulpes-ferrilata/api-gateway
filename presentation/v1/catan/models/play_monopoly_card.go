package models

type PlayMonopolyCardRequest struct {
	DevelopmentCardID         string `json:"developmentCardID" validate:"required,objectid"`
	DemandingResourceCardType string `json:"demandingResourceCardType"`
}
