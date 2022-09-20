package responses

type Game struct {
	ID               string             `json:"id"`
	Status           string             `json:"status"`
	Phase            string             `json:"phase"`
	Turn             int                `json:"turn"`
	Me               *Player            `json:"me"`
	Players          []*Player          `json:"players"`
	Dices            []*Dice            `json:"dices"`
	Achievements     []*Achievement     `json:"achievements"`
	ResourceCards    []*ResourceCard    `json:"resourceCards"`
	DevelopmentCards []*DevelopmentCard `json:"developmentCards"`
	Terrains         []*Terrain         `json:"terrains"`
	Lands            []*Land            `json:"lands"`
	Paths            []*Path            `json:"paths"`
}
