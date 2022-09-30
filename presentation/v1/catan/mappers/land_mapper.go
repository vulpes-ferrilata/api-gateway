package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/responses"
	pb_responses "github.com/vulpes-ferrilata/catan-service-proto/pb/responses"
)

func toLandHttpResponse(landPbResponse *pb_responses.Land) *responses.Land {
	if landPbResponse == nil {
		return nil
	}

	return &responses.Land{
		ID:       landPbResponse.GetID(),
		Q:        int(landPbResponse.GetQ()),
		R:        int(landPbResponse.GetR()),
		Location: landPbResponse.GetLocation(),
	}
}
