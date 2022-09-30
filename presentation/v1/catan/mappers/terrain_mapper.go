package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/responses"
	pb_responses "github.com/vulpes-ferrilata/catan-service-proto/pb/responses"
)

func toTerrainHttpResponse(terrainPbResponse *pb_responses.Terrain) *responses.Terrain {
	if terrainPbResponse == nil {
		return nil
	}

	harborHttpResponses := toHarborHttpResponse(terrainPbResponse.GetHarbor())

	robberHttpResponse := toRobberHttpResponse(terrainPbResponse.GetRobber())

	return &responses.Terrain{
		ID:     terrainPbResponse.GetID(),
		Q:      int(terrainPbResponse.GetQ()),
		R:      int(terrainPbResponse.GetR()),
		Number: int(terrainPbResponse.GetNumber()),
		Type:   terrainPbResponse.GetType(),
		Harbor: harborHttpResponses,
		Robber: robberHttpResponse,
	}
}
