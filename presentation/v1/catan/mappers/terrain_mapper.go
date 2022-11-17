package mappers

import (
	"github.com/pkg/errors"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/responses"
	pb_responses "github.com/vulpes-ferrilata/catan-service-proto/pb/responses"
)

type terrainMapper struct{}

func (t terrainMapper) ToHttpResponse(terrainPbResponse *pb_responses.Terrain) (*responses.Terrain, error) {
	if terrainPbResponse == nil {
		return nil, nil
	}

	harborHttpResponse, err := harborMapper{}.ToHttpResponse(terrainPbResponse.GetHarbor())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	robberHttpResponse, err := robberMapper{}.ToHttpResponse(terrainPbResponse.GetRobber())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &responses.Terrain{
		ID:     terrainPbResponse.GetID(),
		Q:      int(terrainPbResponse.GetQ()),
		R:      int(terrainPbResponse.GetR()),
		Number: int(terrainPbResponse.GetNumber()),
		Type:   terrainPbResponse.GetType(),
		Harbor: harborHttpResponse,
		Robber: robberHttpResponse,
	}, nil
}
