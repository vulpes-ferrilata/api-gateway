package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/utils/slices"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/responses"
	"github.com/vulpes-ferrilata/shared/proto/v1/catan"
)

func ToGameHttpResponse(gameGrpcResponse *catan.GameResponse) *responses.Game {
	if gameGrpcResponse == nil {
		return nil
	}

	me := toPlayerHttpResponse(gameGrpcResponse.Me)

	playerHttpResponses, _ := slices.Map(func(playerGrpcResponse *catan.PlayerResponse) (*responses.Player, error) {
		return toPlayerHttpResponse(playerGrpcResponse), nil
	}, gameGrpcResponse.GetPlayers())

	diceHttpResponses, _ := slices.Map(func(diceGrpcResponse *catan.DiceResponse) (*responses.Dice, error) {
		return toDiceHttpResponse(diceGrpcResponse), nil
	}, gameGrpcResponse.GetDices())

	achievementHttpResponses, _ := slices.Map(func(achievementGrpcResponse *catan.AchievementResponse) (*responses.Achievement, error) {
		return toAchievementHttpResponse(achievementGrpcResponse), nil
	}, gameGrpcResponse.GetAchievements())

	resourceCardHttpResponses, _ := slices.Map(func(resourceCardGrpcResponse *catan.ResourceCardResponse) (*responses.ResourceCard, error) {
		return toResourceCardHttpResponse(resourceCardGrpcResponse), nil
	}, gameGrpcResponse.GetResourceCards())

	developmentCardHttpResponses, _ := slices.Map(func(developmentCardGrpcResponse *catan.DevelopmentCardResponse) (*responses.DevelopmentCard, error) {
		return toDevelopmentCardHttpResponse(developmentCardGrpcResponse), nil
	}, gameGrpcResponse.GetDevelopmentCards())

	terrainHttpResponses, _ := slices.Map(func(terrainGrpcResponse *catan.TerrainResponse) (*responses.Terrain, error) {
		return toTerrainHttpResponse(terrainGrpcResponse), nil
	}, gameGrpcResponse.GetTerrains())

	landHttpResponses, _ := slices.Map(func(landGrpcResponse *catan.LandResponse) (*responses.Land, error) {
		return toLandHttpResponse(landGrpcResponse), nil
	}, gameGrpcResponse.GetLands())

	pathHttpResponses, _ := slices.Map(func(pathGrpcResponse *catan.PathResponse) (*responses.Path, error) {
		return toPathHttpResponse(pathGrpcResponse), nil
	}, gameGrpcResponse.GetPaths())

	return &responses.Game{
		ID:               gameGrpcResponse.GetID(),
		Status:           gameGrpcResponse.GetStatus(),
		Phase:            gameGrpcResponse.GetPhase(),
		Turn:             int(gameGrpcResponse.GetTurn()),
		Me:               me,
		Players:          playerHttpResponses,
		Dices:            diceHttpResponses,
		Achievements:     achievementHttpResponses,
		ResourceCards:    resourceCardHttpResponses,
		DevelopmentCards: developmentCardHttpResponses,
		Terrains:         terrainHttpResponses,
		Lands:            landHttpResponses,
		Paths:            pathHttpResponses,
	}
}
