package requests

type PlayYearOfPlentyCard struct {
	DevelopmentCardID          string   `json:"developmentCardID" validate:"required,objectid"`
	DemandingResourceCardTypes []string `json:"demandingResourceCardTypes" validate:"required,min=1,max=2"`
}
