package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/responses"
	"github.com/vulpes-ferrilata/shared/proto/v1/catan"
)

func toHarborHttpResponse(harborGrpcResponse *catan.HarborResponse) *responses.Harbor {
	if harborGrpcResponse == nil {
		return nil
	}

	return &responses.Harbor{
		ID:   harborGrpcResponse.GetID(),
		Q:    int(harborGrpcResponse.GetQ()),
		R:    int(harborGrpcResponse.GetR()),
		Type: harborGrpcResponse.GetType(),
	}
}
