package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/responses"
	pb_responses "github.com/vulpes-ferrilata/catan-service-proto/pb/responses"
)

type landMapper struct{}

func (l landMapper) ToHttpResponse(landPbResponse *pb_responses.Land) (*responses.Land, error) {
	if landPbResponse == nil {
		return nil, nil
	}

	return &responses.Land{
		ID:       landPbResponse.GetID(),
		Q:        int(landPbResponse.GetQ()),
		R:        int(landPbResponse.GetR()),
		Location: landPbResponse.GetLocation(),
	}, nil
}
