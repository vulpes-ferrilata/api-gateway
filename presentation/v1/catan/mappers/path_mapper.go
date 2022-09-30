package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/responses"
	pb_responses "github.com/vulpes-ferrilata/catan-service-proto/pb/responses"
)

func toPathHttpResponse(pathPbResponse *pb_responses.Path) *responses.Path {
	if pathPbResponse == nil {
		return nil
	}

	return &responses.Path{
		ID:       pathPbResponse.GetID(),
		Q:        int(pathPbResponse.GetQ()),
		R:        int(pathPbResponse.GetR()),
		Location: pathPbResponse.GetLocation(),
	}
}
