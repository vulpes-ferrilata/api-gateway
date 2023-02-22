package models

type SendTradeOfferRequest struct {
	PlayerID string `json:"playerID" validate:"required,objectid"`
}
