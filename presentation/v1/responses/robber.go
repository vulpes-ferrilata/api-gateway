package responses

type Robber struct {
	ID        string `json:"id"`
	TerrainID string `json:"terrainID"`
	IsMoving  bool   `json:"isMoving"`
}
