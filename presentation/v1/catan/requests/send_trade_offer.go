package requests

type SendTradeOffer struct {
	PlayerID string `json:"playerID" validate:"required,objectid"`
}
