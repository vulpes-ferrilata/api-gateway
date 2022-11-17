package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/responses"
	pb_responses "github.com/vulpes-ferrilata/catan-service-proto/pb/responses"
)

type diceMapper struct{}

func (d diceMapper) ToHttpResponse(dicePbResponse *pb_responses.Dice) (*responses.Dice, error) {
	if dicePbResponse == nil {
		return nil, nil
	}

	return &responses.Dice{
		ID:     dicePbResponse.GetID(),
		Number: int(dicePbResponse.GetNumber()),
	}, nil
}
