package requests

type PlayRoadBuildingCard struct {
	PathIDs []string `json:"pathIDs" validate:"required,min=1,max=2,unique"`
}
