package responses

type Harbor struct {
	ID        string `json:"id"`
	TerrainID string `json:"terrainID"`
	Q         int    `json:"q"`
	R         int    `json:"r"`
	Type      string `json:"type"`
}
