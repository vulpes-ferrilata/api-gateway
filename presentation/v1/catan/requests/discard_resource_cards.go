package requests

type DiscardResourceCards struct {
	ResourceCardIDs []string `json:"resourceCardIDs" validate:"required,unique"`
}
