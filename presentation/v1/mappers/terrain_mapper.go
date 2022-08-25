package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/responses"
	"github.com/vulpes-ferrilata/shared/proto/v1/catan"
)

func toTerrainHttpResponse(terrainGrpcResponse *catan.TerrainResponse) *responses.Terrain {
	if terrainGrpcResponse == nil {
		return nil
	}

	return &responses.Terrain{
		ID:     terrainGrpcResponse.GetID(),
		Q:      int(terrainGrpcResponse.GetQ()),
		R:      int(terrainGrpcResponse.GetR()),
		Number: int(terrainGrpcResponse.GetNumber()),
		Type:   terrainGrpcResponse.GetType(),
	}
}
