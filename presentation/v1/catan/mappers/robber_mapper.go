package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/models"
	pb_models "github.com/vulpes-ferrilata/catan-service-proto/pb/models"
)

type robberMapper struct{}

func (r robberMapper) ToHttpResponse(robberPbResponse *pb_models.Robber) (*models.Robber, error) {
	if robberPbResponse == nil {
		return nil, nil
	}

	return &models.Robber{
		ID: robberPbResponse.GetID(),
	}, nil
}
