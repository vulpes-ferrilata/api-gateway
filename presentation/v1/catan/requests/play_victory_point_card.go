package requests

type PlayVictoryPointCard struct {
	DevelopmentCardID string `json:"developmentCardID" validate:"required,objectid"`
}
