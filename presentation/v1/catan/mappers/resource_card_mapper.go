package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/responses"
	pb_responses "github.com/vulpes-ferrilata/catan-service-proto/pb/responses"
)

func toResourceCardHttpResponse(resourceCardPbResponse *pb_responses.ResourceCard) *responses.ResourceCard {
	if resourceCardPbResponse == nil {
		return nil
	}

	return &responses.ResourceCard{
		ID:         resourceCardPbResponse.GetID(),
		Type:       resourceCardPbResponse.GetType(),
		IsSelected: resourceCardPbResponse.GetIsSelected(),
	}
}
