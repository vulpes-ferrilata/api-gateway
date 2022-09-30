package requests

type ToggleResourceCards struct {
	ResourceCardIDs []string `json:"resourceCardIDs" validate:"required,unique"`
}
