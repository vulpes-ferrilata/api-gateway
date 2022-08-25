package responses

import "github.com/vulpes-ferrilata/shared/proto/v1/authentication"

func FromGrpcResponse(tokenGrpcResponse *authentication.TokenResponse) *Token {
	return &Token{
		AccessToken:  tokenGrpcResponse.GetAccessToken(),
		RefreshToken: tokenGrpcResponse.GetRefreshToken(),
	}
}

type Token struct {
	AccessToken  string `json:"accessToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
}
