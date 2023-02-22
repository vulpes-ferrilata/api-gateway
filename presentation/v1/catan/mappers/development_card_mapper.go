package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/models"
	pb_models "github.com/vulpes-ferrilata/catan-service-proto/pb/models"
)

type developmentCardMapper struct{}

func (d developmentCardMapper) ToHttpResponse(developmentCardPbResponse *pb_models.DevelopmentCard) (*models.DevelopmentCard, error) {
	if developmentCardPbResponse == nil {
		return nil, nil
	}

	return &models.DevelopmentCard{
		ID:     developmentCardPbResponse.GetID(),
		Type:   developmentCardPbResponse.GetType(),
		Status: developmentCardPbResponse.GetStatus(),
	}, nil
}
