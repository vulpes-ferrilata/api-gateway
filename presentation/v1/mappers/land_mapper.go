package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/responses"
	"github.com/vulpes-ferrilata/shared/proto/v1/catan"
)

func toLandHttpResponse(landGrpcResponse *catan.LandResponse) *responses.Land {
	if landGrpcResponse == nil {
		return nil
	}

	return &responses.Land{
		ID:       landGrpcResponse.GetID(),
		Q:        int(landGrpcResponse.GetQ()),
		R:        int(landGrpcResponse.GetR()),
		Location: landGrpcResponse.GetLocation(),
	}
}
