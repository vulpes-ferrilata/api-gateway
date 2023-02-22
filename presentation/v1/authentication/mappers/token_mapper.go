package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/authentication/models"
	pb_models "github.com/vulpes-ferrilata/authentication-service-proto/pb/models"
)

func ToTokenHttpResponse(tokenPbResponse *pb_models.Token) *models.Token {
	if tokenPbResponse == nil {
		return nil
	}

	return &models.Token{
		AccessToken:  tokenPbResponse.GetAccessToken(),
		RefreshToken: tokenPbResponse.GetRefreshToken(),
	}
}
