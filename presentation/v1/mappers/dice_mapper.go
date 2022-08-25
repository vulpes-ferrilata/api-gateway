package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/responses"
	"github.com/vulpes-ferrilata/shared/proto/v1/catan"
)

func toDiceHttpResponse(diceGrpcResponse *catan.DiceResponse) *responses.Dice {
	if diceGrpcResponse == nil {
		return nil
	}

	return &responses.Dice{
		ID:     diceGrpcResponse.GetID(),
		Number: int(diceGrpcResponse.GetNumber()),
	}
}
