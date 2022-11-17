package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/responses"
	pb_responses "github.com/vulpes-ferrilata/catan-service-proto/pb/responses"
)

type robberMapper struct{}

func (r robberMapper) ToHttpResponse(robberPbResponse *pb_responses.Robber) (*responses.Robber, error) {
	if robberPbResponse == nil {
		return nil, nil
	}

	return &responses.Robber{
		ID: robberPbResponse.GetID(),
	}, nil
}
