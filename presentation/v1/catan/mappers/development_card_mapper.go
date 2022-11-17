package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/responses"
	pb_responses "github.com/vulpes-ferrilata/catan-service-proto/pb/responses"
)

type developmentCardMapper struct{}

func (d developmentCardMapper) ToHttpResponse(developmentCardPbResponse *pb_responses.DevelopmentCard) (*responses.DevelopmentCard, error) {
	if developmentCardPbResponse == nil {
		return nil, nil
	}

	return &responses.DevelopmentCard{
		ID:     developmentCardPbResponse.GetID(),
		Type:   developmentCardPbResponse.GetType(),
		Status: developmentCardPbResponse.GetStatus(),
	}, nil
}
