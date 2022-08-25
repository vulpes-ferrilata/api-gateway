package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/responses"
	"github.com/vulpes-ferrilata/shared/proto/v1/authentication"
)

func ToTokenHttpResponse(tokenGrpcResponse *authentication.TokenResponse) *responses.Token {
	if tokenGrpcResponse == nil {
		return nil
	}

	return &responses.Token{
		AccessToken:  tokenGrpcResponse.AccessToken,
		RefreshToken: tokenGrpcResponse.RefreshToken,
	}
}
