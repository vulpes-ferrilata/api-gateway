package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/responses"
	pb_responses "github.com/vulpes-ferrilata/catan-service-proto/pb/responses"
)

type resourceCardMapper struct{}

func (r resourceCardMapper) ToHttpResponse(resourceCardPbResponse *pb_responses.ResourceCard) (*responses.ResourceCard, error) {
	if resourceCardPbResponse == nil {
		return nil, nil
	}

	return &responses.ResourceCard{
		ID:       resourceCardPbResponse.GetID(),
		Type:     resourceCardPbResponse.GetType(),
		Offering: resourceCardPbResponse.GetOffering(),
	}, nil
}
