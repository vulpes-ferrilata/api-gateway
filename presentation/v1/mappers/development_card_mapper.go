package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/responses"
	"github.com/vulpes-ferrilata/shared/proto/v1/catan"
)

func toDevelopmentCardHttpResponse(developmentCardGrpcResponse *catan.DevelopmentCardResponse) *responses.DevelopmentCard {
	if developmentCardGrpcResponse == nil {
		return nil
	}

	return &responses.DevelopmentCard{
		ID:     developmentCardGrpcResponse.GetID(),
		Type:   developmentCardGrpcResponse.GetType(),
		Status: developmentCardGrpcResponse.GetStatus(),
	}
}
