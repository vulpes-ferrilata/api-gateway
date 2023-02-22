package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/models"
	pb_models "github.com/vulpes-ferrilata/catan-service-proto/pb/models"
)

type resourceCardMapper struct{}

func (r resourceCardMapper) ToHttpResponse(resourceCardPbResponse *pb_models.ResourceCard) (*models.ResourceCard, error) {
	if resourceCardPbResponse == nil {
		return nil, nil
	}

	return &models.ResourceCard{
		ID:       resourceCardPbResponse.GetID(),
		Type:     resourceCardPbResponse.GetType(),
		Offering: resourceCardPbResponse.GetOffering(),
	}, nil
}
