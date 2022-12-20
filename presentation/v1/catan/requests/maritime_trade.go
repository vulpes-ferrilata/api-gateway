package requests

type MaritimeTrade struct {
	ResourceCardType          string `json:"resourceCardType"`
	DemandingResourceCardType string `json:"demandingResourceCardType"`
}
