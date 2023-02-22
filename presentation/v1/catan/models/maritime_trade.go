package models

type MaritimeTradeRequest struct {
	ResourceCardType          string `json:"resourceCardType"`
	DemandingResourceCardType string `json:"demandingResourceCardType"`
}
