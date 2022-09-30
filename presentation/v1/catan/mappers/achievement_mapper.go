package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/responses"
	pb_responses "github.com/vulpes-ferrilata/catan-service-proto/pb/responses"
)

func toAchievementHttpResponse(achievementPbResponse *pb_responses.Achievement) *responses.Achievement {
	if achievementPbResponse == nil {
		return nil
	}

	return &responses.Achievement{
		ID:   achievementPbResponse.GetID(),
		Type: achievementPbResponse.GetType(),
	}
}
