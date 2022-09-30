package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/responses"
	pb_responses "github.com/vulpes-ferrilata/catan-service-proto/pb/responses"
)

func toRobberHttpResponse(robberPbResponse *pb_responses.Robber) *responses.Robber {
	if robberPbResponse == nil {
		return nil
	}

	return &responses.Robber{
		ID: robberPbResponse.GetID(),
	}
}
