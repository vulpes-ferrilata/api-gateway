package requests

type PlayRoadBuildingCard struct {
	DevelopmentCardID string   `json:"developmentCardID" validate:"required,objectid"`
	PathIDs           []string `json:"pathIDs" validate:"required,min=1,max=2,unique"`
}
