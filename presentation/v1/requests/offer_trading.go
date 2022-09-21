package requests

type OfferTrading struct {
	PlayerID string `json:"playerID" validate:"required,objectid"`
}
