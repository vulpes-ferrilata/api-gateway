package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/responses"
	"github.com/vulpes-ferrilata/shared/proto/v1/catan"
)

func toPathHttpResponse(pathGrpcResponse *catan.PathResponse) *responses.Path {
	if pathGrpcResponse == nil {
		return nil
	}

	return &responses.Path{
		ID:       pathGrpcResponse.GetID(),
		Q:        int(pathGrpcResponse.GetQ()),
		R:        int(pathGrpcResponse.GetR()),
		Location: pathGrpcResponse.GetLocation(),
	}
}
