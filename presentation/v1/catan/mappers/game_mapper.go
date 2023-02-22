package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/models"
	pb_models "github.com/vulpes-ferrilata/catan-service-proto/pb/models"
)

type gameMapper struct{}

func (g gameMapper) ToHttpResponse(gamePbResponse *pb_models.Game) (*models.Game, error) {
	if gamePbResponse == nil {
		return nil, nil
	}

	return &models.Game{
		ID:             gamePbResponse.GetID(),
		PlayerQuantity: int(gamePbResponse.GetPlayerQuantity()),
		Status:         gamePbResponse.GetStatus(),
	}, nil
}
