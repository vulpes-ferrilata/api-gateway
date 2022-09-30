package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/authentication/responses"
	pb_responses "github.com/vulpes-ferrilata/authentication-service-proto/pb/responses"
)

func ToTokenHttpResponse(tokenPbResponse *pb_responses.Token) *responses.Token {
	if tokenPbResponse == nil {
		return nil
	}

	return &responses.Token{
		AccessToken:  tokenPbResponse.GetAccessToken(),
		RefreshToken: tokenPbResponse.GetRefreshToken(),
	}
}
