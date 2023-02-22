package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/models"
	pb_models "github.com/vulpes-ferrilata/catan-service-proto/pb/models"
)

type harborMapper struct{}

func (h harborMapper) ToHttpResponse(harborPbResponse *pb_models.Harbor) (*models.Harbor, error) {
	if harborPbResponse == nil {
		return nil, nil
	}

	return &models.Harbor{
		ID:   harborPbResponse.GetID(),
		Q:    int(harborPbResponse.GetQ()),
		R:    int(harborPbResponse.GetR()),
		Type: harborPbResponse.GetType(),
	}, nil
}
