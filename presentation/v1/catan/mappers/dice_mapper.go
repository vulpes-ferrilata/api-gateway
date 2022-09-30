package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/responses"
	pb_responses "github.com/vulpes-ferrilata/catan-service-proto/pb/responses"
)

func toDiceHttpResponse(dicePbResponse *pb_responses.Dice) *responses.Dice {
	if dicePbResponse == nil {
		return nil
	}

	return &responses.Dice{
		ID:     dicePbResponse.GetID(),
		Number: int(dicePbResponse.GetNumber()),
	}
}
