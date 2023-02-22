package models

type BuildSettlementAndRoadRequest struct {
	LandID string `json:"landID" validate:"required,objectid"`
	PathID string `json:"pathID" validate:"required,objectid"`
}
