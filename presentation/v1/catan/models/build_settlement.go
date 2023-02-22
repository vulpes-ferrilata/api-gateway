package models

type BuildSettlementRequest struct {
	LandID string `json:"landID" validate:"required,objectid"`
}
