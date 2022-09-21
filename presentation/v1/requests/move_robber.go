package requests

type MoveRobber struct {
	TerrainID string `json:"terrainID" validate:"required,objectid"`
	PlayerID  string `json:"playerID" validate:"required,objectid"`
}
