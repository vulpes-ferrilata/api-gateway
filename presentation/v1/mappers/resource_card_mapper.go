package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/responses"
	"github.com/vulpes-ferrilata/shared/proto/v1/catan"
)

func toResourceCardHttpResponse(resourceCardGrpcResponse *catan.ResourceCardResponse) *responses.ResourceCard {
	if resourceCardGrpcResponse == nil {
		return nil
	}

	return &responses.ResourceCard{
		ID:         resourceCardGrpcResponse.GetID(),
		Type:       resourceCardGrpcResponse.GetType(),
		IsSelected: resourceCardGrpcResponse.GetIsSelected(),
	}
}
