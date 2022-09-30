package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/responses"
	pb_responses "github.com/vulpes-ferrilata/catan-service-proto/pb/responses"
)

func toRoadHttpResponse(roadPbResponse *pb_responses.Road) *responses.Road {
	if roadPbResponse == nil {
		return nil
	}

	pathHttpResponse := toPathHttpResponse(roadPbResponse.GetPath())

	return &responses.Road{
		ID:   roadPbResponse.GetID(),
		Path: pathHttpResponse,
	}
}
