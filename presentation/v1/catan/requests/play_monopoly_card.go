package requests

type PlayMonopolyCard struct {
	DevelopmentCardID         string `json:"developmentCardID" validate:"required,objectid"`
	DemandingResourceCardType string `json:"demandingResourceCardType"`
}
