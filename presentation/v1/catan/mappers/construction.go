package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/responses"
	pb_responses "github.com/vulpes-ferrilata/catan-service-proto/pb/responses"
)

func toConstructionHttpResponse(constructionPbResponse *pb_responses.Construction) *responses.Construction {
	if constructionPbResponse == nil {
		return nil
	}

	landHttpResponse := toLandHttpResponse(constructionPbResponse.GetLand())

	return &responses.Construction{
		ID:   constructionPbResponse.GetID(),
		Type: constructionPbResponse.GetType(),
		Land: landHttpResponse,
	}
}
