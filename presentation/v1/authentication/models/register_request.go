package models

type RegisterRequest struct {
	DisplayName string `json:"displayName" validate:"required,min=1,max=20"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
}
