package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/models"
	pb_models "github.com/vulpes-ferrilata/catan-service-proto/pb/models"
)

type diceMapper struct{}

func (d diceMapper) ToHttpResponse(dicePbResponse *pb_models.Dice) (*models.Dice, error) {
	if dicePbResponse == nil {
		return nil, nil
	}

	return &models.Dice{
		ID:     dicePbResponse.GetID(),
		Number: int(dicePbResponse.GetNumber()),
	}, nil
}
