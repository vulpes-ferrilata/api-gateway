package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/responses"
	"github.com/vulpes-ferrilata/shared/proto/v1/catan"
)

func toRobberHttpResponse(robberGrpcResponse *catan.RobberResponse) *responses.Robber {
	if robberGrpcResponse == nil {
		return nil
	}

	return &responses.Robber{
		ID: robberGrpcResponse.GetID(),
	}
}
