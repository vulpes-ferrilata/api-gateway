package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/models"
	pb_models "github.com/vulpes-ferrilata/catan-service-proto/pb/models"
)

type achievementMapper struct{}

func (a achievementMapper) ToHttpResponse(achievementPbResponse *pb_models.Achievement) (*models.Achievement, error) {
	if achievementPbResponse == nil {
		return nil, nil
	}

	return &models.Achievement{
		ID:   achievementPbResponse.GetID(),
		Type: achievementPbResponse.GetType(),
	}, nil
}
