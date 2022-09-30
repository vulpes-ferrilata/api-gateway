package requests

type BuildSettlement struct {
	LandID string `json:"landID" validate:"required,objectid"`
}
