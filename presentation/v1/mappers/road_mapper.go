package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/responses"
	"github.com/vulpes-ferrilata/shared/proto/v1/catan"
)

func toRoadHttpResponse(roadGrpcResponse *catan.RoadResponse) *responses.Road {
	if roadGrpcResponse == nil {
		return nil
	}

	pathHttpResponse := toPathHttpResponse(roadGrpcResponse.GetPath())

	return &responses.Road{
		ID:   roadGrpcResponse.GetID(),
		Path: pathHttpResponse,
	}
}
