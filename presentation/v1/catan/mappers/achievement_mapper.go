package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/responses"
	pb_responses "github.com/vulpes-ferrilata/catan-service-proto/pb/responses"
)

type achievementMapper struct{}

func (a achievementMapper) ToHttpResponse(achievementPbResponse *pb_responses.Achievement) (*responses.Achievement, error) {
	if achievementPbResponse == nil {
		return nil, nil
	}

	return &responses.Achievement{
		ID:   achievementPbResponse.GetID(),
		Type: achievementPbResponse.GetType(),
	}, nil
}
