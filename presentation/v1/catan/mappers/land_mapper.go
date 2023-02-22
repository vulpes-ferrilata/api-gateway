package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/models"
	pb_models "github.com/vulpes-ferrilata/catan-service-proto/pb/models"
)

type landMapper struct{}

func (l landMapper) ToHttpResponse(landPbResponse *pb_models.Land) (*models.Land, error) {
	if landPbResponse == nil {
		return nil, nil
	}

	return &models.Land{
		ID:       landPbResponse.GetID(),
		Q:        int(landPbResponse.GetQ()),
		R:        int(landPbResponse.GetR()),
		Location: landPbResponse.GetLocation(),
	}, nil
}
