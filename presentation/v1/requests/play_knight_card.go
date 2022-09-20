package requests

type PlayKnightCard struct {
	TerrainID string `json:"terrainID"`
	PlayerID  string `json:"playerID"`
}
