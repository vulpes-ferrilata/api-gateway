package responses

type Game struct {
	ID               string             `json:"id"`
	Status           string             `json:"status"`
	Turn             int                `json:"turn"`
	IsRolledDices    bool               `json:"isRolledDices"`
	Players          []*Player          `json:"players"`
	Dices            []*Dice            `json:"dices"`
	Achievements     []*Achievement     `json:"achievements"`
	ResourceCards    []*ResourceCard    `json:"resourceCards"`
	DevelopmentCards []*DevelopmentCard `json:"developmentCards"`
	Terrains         []*Terrain         `json:"terrains"`
	Harbors          []*Harbor          `json:"harbors"`
	Robber           *Robber            `json:"robber"`
	Lands            []*Land            `json:"lands"`
	Paths            []*Path            `json:"paths"`
}
