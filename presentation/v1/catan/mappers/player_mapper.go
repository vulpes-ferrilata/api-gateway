package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/utils/slices"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/responses"
	pb_responses "github.com/vulpes-ferrilata/catan-service-proto/pb/responses"
)

func toPlayerHttpResponse(playerPbResponse *pb_responses.Player) *responses.Player {
	if playerPbResponse == nil {
		return nil
	}

	achievementHttpResponses, _ := slices.Map(func(achievementPbResponse *pb_responses.Achievement) (*responses.Achievement, error) {
		return toAchievementHttpResponse(achievementPbResponse), nil
	}, playerPbResponse.GetAchievements())

	resourceCardHttpResponses, _ := slices.Map(func(resourceCardPbResponse *pb_responses.ResourceCard) (*responses.ResourceCard, error) {
		return toResourceCardHttpResponse(resourceCardPbResponse), nil
	}, playerPbResponse.GetResourceCards())

	developmentCardHttpResponses, _ := slices.Map(func(developmentCardPbResponse *pb_responses.DevelopmentCard) (*responses.DevelopmentCard, error) {
		return toDevelopmentCardHttpResponse(developmentCardPbResponse), nil
	}, playerPbResponse.GetDevelopmentCards())

	constructionHttpResponses, _ := slices.Map(func(constructionPbResponse *pb_responses.Construction) (*responses.Construction, error) {
		return toConstructionHttpResponse(constructionPbResponse), nil
	}, playerPbResponse.GetConstructions())

	roadHttpResponses, _ := slices.Map(func(roadPbResponse *pb_responses.Road) (*responses.Road, error) {
		return toRoadHttpResponse(roadPbResponse), nil
	}, playerPbResponse.GetRoads())

	return &responses.Player{
		ID:               playerPbResponse.GetID(),
		UserID:           playerPbResponse.GetUserID(),
		Color:            playerPbResponse.GetColor(),
		TurnOrder:        int(playerPbResponse.GetTurnOrder()),
		IsOffered:        playerPbResponse.GetIsOffered(),
		Score:            int(playerPbResponse.GetScore()),
		Achievements:     achievementHttpResponses,
		ResourceCards:    resourceCardHttpResponses,
		DevelopmentCards: developmentCardHttpResponses,
		Constructions:    constructionHttpResponses,
		Roads:            roadHttpResponses,
	}
}
