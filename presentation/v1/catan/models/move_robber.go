package models

type MoveRobberRequest struct {
	TerrainID string `json:"terrainID" validate:"required,objectid"`
	PlayerID  string `json:"playerID" validate:"omitempty,objectid"`
}
