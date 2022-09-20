package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/utils/slices"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/responses"
	"github.com/vulpes-ferrilata/shared/proto/v1/catan"
)

func toPlayerHttpResponse(playerGrpcResponse *catan.PlayerResponse) *responses.Player {
	if playerGrpcResponse == nil {
		return nil
	}

	achievementHttpResponses, _ := slices.Map(func(achievementGrpcResponse *catan.AchievementResponse) (*responses.Achievement, error) {
		return toAchievementHttpResponse(achievementGrpcResponse), nil
	}, playerGrpcResponse.GetAchievements())

	resourceCardHttpResponses, _ := slices.Map(func(resourceCardGrpcResponse *catan.ResourceCardResponse) (*responses.ResourceCard, error) {
		return toResourceCardHttpResponse(resourceCardGrpcResponse), nil
	}, playerGrpcResponse.GetResourceCards())

	developmentCardHttpResponses, _ := slices.Map(func(developmentCardGrpcResponse *catan.DevelopmentCardResponse) (*responses.DevelopmentCard, error) {
		return toDevelopmentCardHttpResponse(developmentCardGrpcResponse), nil
	}, playerGrpcResponse.GetDevelopmentCards())

	constructionHttpResponses, _ := slices.Map(func(constructionGrpcResponse *catan.ConstructionResponse) (*responses.Construction, error) {
		return toConstructionHttpResponse(constructionGrpcResponse), nil
	}, playerGrpcResponse.GetConstructions())

	roadHttpResponses, _ := slices.Map(func(roadGrpcResponse *catan.RoadResponse) (*responses.Road, error) {
		return toRoadHttpResponse(roadGrpcResponse), nil
	}, playerGrpcResponse.GetRoads())

	return &responses.Player{
		ID:               playerGrpcResponse.GetID(),
		UserID:           playerGrpcResponse.GetUserID(),
		Color:            playerGrpcResponse.GetColor(),
		TurnOrder:        int(playerGrpcResponse.GetTurnOrder()),
		IsOffered:        playerGrpcResponse.GetIsOffered(),
		IsActive:         playerGrpcResponse.GetIsActive(),
		Score:            int(playerGrpcResponse.GetScore()),
		Achievements:     achievementHttpResponses,
		ResourceCards:    resourceCardHttpResponses,
		DevelopmentCards: developmentCardHttpResponses,
		Constructions:    constructionHttpResponses,
		Roads:            roadHttpResponses,
	}
}
