package models

type BuildRoadRequest struct {
	PathID string `json:"pathID" validate:"required,objectid"`
}
