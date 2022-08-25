package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/responses"
	"github.com/vulpes-ferrilata/shared/proto/v1/catan"
)

func toAchievementHttpResponse(achievementGrpcResponse *catan.AchievementResponse) *responses.Achievement {
	if achievementGrpcResponse == nil {
		return nil
	}

	return &responses.Achievement{
		ID:   achievementGrpcResponse.GetID(),
		Type: achievementGrpcResponse.GetType(),
	}
}
