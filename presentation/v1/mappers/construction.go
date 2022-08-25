package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/responses"
	"github.com/vulpes-ferrilata/shared/proto/v1/catan"
)

func toConstructionHttpResponse(constructionGrpcResponse *catan.ConstructionResponse) *responses.Construction {
	if constructionGrpcResponse == nil {
		return nil
	}

	landHttpResponse := toLandHttpResponse(constructionGrpcResponse.GetLand())

	return &responses.Construction{
		ID:   constructionGrpcResponse.GetID(),
		Type: constructionGrpcResponse.GetType(),
		Land: landHttpResponse,
	}
}
