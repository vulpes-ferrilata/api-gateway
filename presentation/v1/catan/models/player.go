package models

type Player struct {
	ID                 string             `json:"id"`
	UserID             string             `json:"userID"`
	Color              string             `json:"color"`
	TurnOrder          int                `json:"turnOrder"`
	ReceivedOffer      bool               `json:"receivedOffer"`
	DiscardedResources bool               `json:"discardedResources"`
	Score              int                `json:"score"`
	Achievements       []*Achievement     `json:"achievements"`
	ResourceCards      []*ResourceCard    `json:"resourceCards"`
	DevelopmentCards   []*DevelopmentCard `json:"developmentCards"`
	Constructions      []*Construction    `json:"constructions"`
	Roads              []*Road            `json:"roads"`
}
