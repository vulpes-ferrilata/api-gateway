package requests

type PlayYearOfPlentyCard struct {
	ResourceCardTypes []string `json:"resourceCardTypes" validate:"required,min=1,max=2"`
}
