package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/responses"
	pb_responses "github.com/vulpes-ferrilata/catan-service-proto/pb/responses"
)

type gameMapper struct{}

func (g gameMapper) ToHttpResponse(gamePbResponse *pb_responses.Game) (*responses.Game, error) {
	if gamePbResponse == nil {
		return nil, nil
	}

	return &responses.Game{
		ID:             gamePbResponse.GetID(),
		PlayerQuantity: int(gamePbResponse.GetPlayerQuantity()),
		Status:         gamePbResponse.GetStatus(),
	}, nil
}
