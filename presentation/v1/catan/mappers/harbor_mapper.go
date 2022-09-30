package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/responses"
	pb_responses "github.com/vulpes-ferrilata/catan-service-proto/pb/responses"
)

func toHarborHttpResponse(harborPbResponse *pb_responses.Harbor) *responses.Harbor {
	if harborPbResponse == nil {
		return nil
	}

	return &responses.Harbor{
		ID:   harborPbResponse.GetID(),
		Q:    int(harborPbResponse.GetQ()),
		R:    int(harborPbResponse.GetR()),
		Type: harborPbResponse.GetType(),
	}
}
